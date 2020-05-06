package user

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/casbin/casbin/v2"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/internal/utils"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type User struct {
	oldmodels.User
}

type Service interface {
	Count() (uint64, error)
	Query(limit, offset uint64, qps *QueryParams) ([]oldmodels.User, error)
	Create(form *CreateForm) (*models.User, error)
}

func NewService(db *sqlex.DB, enforcer *casbin.Enforcer) Service {
	return &service{db: db, enforcer: enforcer}
}

type service struct {
	db       *sqlex.DB
	enforcer *casbin.Enforcer
}

func (s *service) Count() (count uint64, err error) {
	sql, args, err := squirrel.Select("count(*)").From("users").ToSql()
	if err != nil {
		return 0, err
	}

	err = s.db.Get(&count, sql, args...)
	return
}

func (s *service) Query(limit, offset uint64, qps *QueryParams) (users []oldmodels.User, err error) {
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

	users = []oldmodels.User{}
	err = s.db.Select(&users, sql, args...)
	return
}

func (s *service) Create(form *CreateForm) (u *models.User, err error) {
	ctx := context.TODO()
	err = utils.Tx(ctx, func(tx *sql.Tx) (err error) {
		u = &models.User{}
		u.Email = form.Email
		u.Username = form.Username
		u.HashedPassword = form.Password
		u.State = form.State
		err = u.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return
		}

		userID := "user_" + strconv.FormatInt(u.ID, 10)
		for _, role := range form.Roles {
			_, err = s.enforcer.AddRoleForUser(userID, role)
			if err != nil {
				return
			}
		}

		return nil
	})
	return
}
