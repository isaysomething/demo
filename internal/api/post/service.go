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

type Service interface {
	Get(id int64) (*models.Post, error)
	Create(*clevergo.Context) (*models.Post, error)
	Count() (int, error)
	Query(limit, offset int, qps *QueryParams) ([]models.Post, error)
	Update(id int64, form *Form) (*models.Post, error)
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

func (s *service) Get(id int64) (post *models.Post, err error) {
	sql, args, err := squirrel.Select("*").From("posts").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	post = new(models.Post)
	err = s.db.Get(post, sql, args...)
	return
}

func (s *service) Create(ctx *clevergo.Context) (post *models.Post, err error) {
	user, _ := s.userManager.Get(ctx.Request, ctx.Response)
	userID, _ := strconv.ParseInt(user.GetIdentity().GetID(), 10, 64)
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return nil, err
	}
	now := time.Now()
	post = &models.Post{
		Title:     form.Title,
		UserID:    userID,
		Content:   form.Content,
		State:     form.State,
		CreatedAt: now,
		UpdatedAt: sqlex.ToNullTime(now),
	}
	err = post.Save(s.db)
	return
}

func (s *service) Update(id int64, form *Form) (post *models.Post, err error) {
	post, err = s.Get(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	post.Title = form.Title
	post.Content = form.Content
	post.State = form.State
	post.UpdatedAt = sqlex.ToNullTime(now)
	err = post.Update(s.db)
	return
}

func (s *service) Count() (count int, err error) {
	sql, args, err := squirrel.Select("COUNT(*)").From("posts").ToSql()
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
