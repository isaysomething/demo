package post

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Post struct {
	models.Post
	MarkdownContent string `json:"markdown_content"`
}

type Service interface {
	Get(id int64) (*Post, error)
	Create(*clevergo.Context) (*Post, error)
	Count() (int, error)
	Query(limit, offset int, qps *QueryParams) ([]models.Post, error)
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

func (s *service) Get(id int64) (post *Post, err error) {
	sql, args, err := squirrel.Select("*").
		From("posts").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	post = new(Post)
	err = s.db.Get(post, sql, args...)
	if err != nil {
		return
	}

	markdownContent, err := models.GetPostMeta(s.db, id, models.PostMetaMarkdownContent)
	if err != nil {
		return nil, err
	}
	post.MarkdownContent = markdownContent.Value

	return
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
	tx, err := s.db.Begin()
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
	sql, args, err = squirrel.Insert("post_meta").SetMap(clevergo.Map{
		"post_id":    postID,
		"meta_key":   models.PostMetaMarkdownContent,
		"meta_value": form.MarkdownContent,
		"created_at": now,
		"updated_at": sqlex.ToNullTime(now),
	}).ToSql()
	if _, err = tx.Exec(sql, args...); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return s.Get(postID)
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
	sql, args, err = squirrel.Update("post_meta").
		Where(squirrel.Eq{
			"post_id":  id,
			"meta_key": models.PostMetaMarkdownContent,
		}).
		SetMap(clevergo.Map{
			"meta_value": form.MarkdownContent,
			"updated_at": sqlex.ToNullTime(now),
		}).
		ToSql()
	if _, err = tx.Exec(sql, args...); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return
}

func (s *service) Count() (count int, err error) {
	sql, args, err := squirrel.Select("COUNT(*)").From("posts").
		Where(squirrel.NotEq{"state": models.PostStateDeleted}).
		ToSql()
	if err != nil {
		return 0, err
	}
	err = s.db.Get(&count, sql, args...)
	return
}

var states = map[string]int{
	"published": models.PostStatePublished,
	"draft":     models.PostStateDraft,
}

func (s *service) Query(limit, offset int, qps *QueryParams) ([]models.Post, error) {
	query := squirrel.Select("*").From("posts").Where(squirrel.NotEq{"state": models.PostStateDeleted})
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

	posts := []models.Post{}
	err = s.db.Select(&posts, sql, args...)
	return posts, err
}

func (s *service) Delete(id int64) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}
	sql, args, err := squirrel.Update("posts").
		Set("state", models.PostStateDeleted).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sql, args...)
	return err
}
