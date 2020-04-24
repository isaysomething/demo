package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Service interface {
	Count() (int, error)
	Query(limit, offset int) ([]models.User, error)
}

func NewService(db *sqlex.DB) Service {
	return &service{db: db}
}

type service struct {
	db *sqlex.DB
}

func (s *service) Count() (count int, err error) {
	sql, args, err := squirrel.Select("count(*)").From("users").ToSql()
	if err != nil {
		return 0, err
	}

	err = s.db.Get(&count, sql, args...)
	return
}

func (s *service) Query(limit, offset int) (users []models.User, err error) {
	sql, args, err := squirrel.Select("*").From("users").ToSql()
	if err != nil {
		return nil, err
	}

	users = []models.User{}
	err = s.db.Select(&users, sql, args...)
	return
}
