package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/sqlex"
)

type User struct {
	models.User
}

type Service interface {
	Count() (uint64, error)
	Query(limit, offset uint64, qps *QueryParams) ([]models.User, error)
	Create(form *CreateForm) (*User, error)
}

func NewService(db *sqlex.DB) Service {
	return &service{db: db}
}

type service struct {
	db *sqlex.DB
}

func (s *service) Count() (count uint64, err error) {
	sql, args, err := squirrel.Select("count(*)").From("users").ToSql()
	if err != nil {
		return 0, err
	}

	err = s.db.Get(&count, sql, args...)
	return
}

func (s *service) Query(limit, offset uint64, qps *QueryParams) (users []models.User, err error) {
	query := squirrel.Select("*").From("users")
	if qps.Username != "" {
		query = query.Where(squirrel.Like{"username": "%" + qps.Username + "%"})
	}
	if qps.Email != "" {
		query = query.Where(squirrel.Like{"email": "%" + qps.Email + "%"})
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

	users = []models.User{}
	err = s.db.Select(&users, sql, args...)
	return
}

func (s *service) Create(form *CreateForm) (*User, error) {
	user, err := models.CreateUser(s.db, form.Username, form.Email, form.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		User: *user,
	}, nil
}
