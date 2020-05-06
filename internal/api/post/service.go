package post

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Post struct {
	oldmodels.Post
	MarkdownContent string `json:"markdown_content"`
}

type Service interface {
	Get(id int64) (*Post, error)
	Create(*clevergo.Context) (*Post, error)
	Count() (uint64, error)
	Query(limit, offset uint64, qps *QueryParams) ([]oldmodels.Post, error)
	Update(id int64, form *Form) (*Post, error)
	Delete(id int64) error
}

type service struct {
	db          *sqlex.DB
	userManager *api.UserManager
}

func NewService(db *sqlex.DB, userManager *api.UserManager) Service {
	return &service{
		db:          db,
		userManager: userManager,
	}
}

func (s *service) Get(id int64) (*Post, error) {
	sql, args, err := squirrel.Select("*").
		From("posts").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	post := new(Post)
	err = s.db.Get(post, sql, args...)
	if err != nil {
		return nil, err
	}

	markdownContent, err := oldmodels.GetPostMeta(s.db, id, oldmodels.PostMetaMarkdownContent)
	if err == nil {
		post.MarkdownContent = markdownContent.Value
	}
	return post, nil
}

func (s *service) Create(ctx *clevergo.Context) (*Post, error) {
	user, _ := s.userManager.Get(ctx.Request, ctx.Response)
	userID, _ := strconv.ParseInt(user.GetIdentity().GetID(), 10, 64)
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return nil, err
	}
	if err := form.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	post, err := s.createPost(tx, userID, form, now)
	if err != nil {
		return nil, err
	}
	if err = s.createPostMeta(tx, post.ID, oldmodels.PostMetaMarkdownContent, form.MarkdownContent, now); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *service) createPost(tx *sql.Tx, userID int64, form *Form, now time.Time) (*Post, error) {
	sql, args, err := squirrel.Insert("posts").
		SetMap(clevergo.Map{
			"user_id":    userID,
			"title":      form.Title,
			"content":    form.Content,
			"state":      form.State,
			"created_at": now,
			"updated_at": sqlex.ToNullTime(now),
		}).ToSql()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(sql, args...)
	if err != nil {
		return nil, err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.Get(postID)
}

func (s *service) createPostMeta(tx *sql.Tx, postID int64, key, value string, now time.Time) error {
	sql, args, err := squirrel.Insert("post_meta").SetMap(clevergo.Map{
		"post_id":    postID,
		"meta_key":   key,
		"meta_value": value,
		"created_at": now,
		"updated_at": sqlex.ToNullTime(now),
	}).ToSql()
	if _, err = tx.Exec(sql, args...); err != nil {
		return err
	}
	return nil
}

func (s *service) updatePostMeta(tx *sql.Tx, postID int64, key, value string, now time.Time) error {
	_, err := oldmodels.GetPostMeta(s.db, postID, key)
	if err != nil && sql.ErrNoRows == err {
		return s.createPostMeta(tx, postID, key, value, now)
	}
	sql, args, err := squirrel.Update("post_meta").
		Where(squirrel.Eq{
			"post_id":  postID,
			"meta_key": key,
		}).
		SetMap(clevergo.Map{
			"meta_value": value,
			"updated_at": sqlex.ToNullTime(now),
		}).
		ToSql()
	if _, err = tx.Exec(sql, args...); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(id int64, form *Form) (post *Post, err error) {
	post, err = s.Get(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sql, args, err := squirrel.Update("posts").
		Where(squirrel.Eq{"id": id}).
		SetMap(clevergo.Map{
			"title":      form.Title,
			"content":    form.Content,
			"state":      form.State,
			"updated_at": sqlex.ToNullTime(now),
		}).ToSql()
	if err != nil {
		return nil, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	if _, err = tx.Exec(sql, args...); err != nil {
		return nil, err
	}
	if err = s.updatePostMeta(tx, id, oldmodels.PostMetaMarkdownContent, form.MarkdownContent, now); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return
}

func (s *service) Count() (count uint64, err error) {
	sql, args, err := squirrel.Select("COUNT(*)").From("posts").
		Where(squirrel.NotEq{"state": oldmodels.PostStateDeleted}).
		ToSql()
	if err != nil {
		return 0, err
	}
	err = s.db.Get(&count, sql, args...)
	return
}

var states = map[string]int{
	"published": oldmodels.PostStatePublished,
	"draft":     oldmodels.PostStateDraft,
}

func (s *service) Query(limit, offset uint64, qps *QueryParams) ([]oldmodels.Post, error) {
	query := squirrel.Select("*").From("posts").Where(squirrel.NotEq{"state": oldmodels.PostStateDeleted})
	if qps.Title != "" {
		query = query.Where(squirrel.Like{"title": fmt.Sprintf("%%%s%%", qps.Title)})
	}
	if qps.State != "" {
		query = query.Where(squirrel.Eq{"state": qps.State})
	}
	if orderBy := qps.OrderBy(); orderBy != "" {
		query = query.OrderBy(orderBy)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	posts := []oldmodels.Post{}
	err = s.db.Select(&posts, sql, args...)
	return posts, err
}

func (s *service) Delete(id int64) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}
	sql, args, err := squirrel.Update("posts").
		Set("state", oldmodels.PostStateDeleted).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sql, args...)
	return err
}
