package authz

import (
	"github.com/Masterminds/squirrel"
	"github.com/clevergo/demo/internal/rbac"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Service interface {
	Count() (uint64, error)
	Query(limit, offset uint64) ([]rbac.Item, error)
}

func NewService(db *sqlex.DB) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db *sqlex.DB
}

func (s *service) Count() (count uint64, err error) {
	sql, args, err := squirrel.Select("count(*)").From("auth_items").
		Where(squirrel.Eq{"item_type": rbac.TypeRole}).
		ToSql()
	if err != nil {
		return
	}
	err = s.db.Get(&count, sql, args...)
	return
}

func (s *service) Query(limit, offset uint64) (roles []rbac.Item, err error) {
	sql, args, err := squirrel.Select("*").From("auth_items").
		Where(squirrel.Eq{"item_type": rbac.TypeRole}).
		OrderBy("id ASC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return
	}
	roles = []rbac.Item{}
	err = s.db.Select(&roles, sql, args...)
	return
}
