package post

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Service interface {
	Get(id int64) (*models.Post, error)
	Create(form *Form) (*models.Post, error)
	Count() (int, error)
	Query(limit, offset int) ([]models.Post, error)
	Update(id int64, form *Form) (*models.Post, error)
	Delete(id int64) error
}

type service struct {
	db *sqlex.DB
}

func NewService(db *sqlex.DB) Service {
	return &service{
		db: db,
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

func (s *service) Create(form *Form) (post *models.Post, err error) {
	now := time.Now()
	post = &models.Post{
		Title:     form.Title,
		Content:   form.Content,
		Status:    form.Status,
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
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
	post.Status = form.Status
	post.UpdatedAt = sql.NullTime{Time: now, Valid: true}
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

func (s *service) Query(limit, offset int) ([]models.Post, error) {
	sql, args, err := squirrel.Select("*").From("posts").ToSql()
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
	sql, args, err := squirrel.Delete("posts").Where(squirrel.Eq{
		"id": id,
	}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sql, args...)
	return err
}
