package authz

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/rbac"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Role struct {
	rbac.Item
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}

type Service interface {
	Count() (uint64, error)
	Query(limit, offset uint64, qps *QueryParams) ([]rbac.Item, error)
	Get(id string) (*Role, error)
	Create(form *Form) (*Role, error)
	Delete(id string) error
}

func NewService(db *sqlex.DB, enforcer *casbin.Enforcer) Service {
	return &service{
		db:       db,
		enforcer: enforcer,
	}
}

type service struct {
	enforcer *casbin.Enforcer
	db       *sqlex.DB
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

func (s *service) Query(limit, offset uint64, qps *QueryParams) (roles []rbac.Item, err error) {
	query := squirrel.Select("*").From("auth_items").
		Where(squirrel.Eq{"item_type": rbac.TypeRole}).
		OrderBy("id ASC").
		Limit(limit).
		Offset(offset)
	if qps.Name != "" {
		query = query.Where(squirrel.Like{"name": "%" + qps.Name + "%"})
	}
	if len(qps.Exclude) > 0 {
		fmt.Println(qps.Exclude)
		query = query.Where(squirrel.NotEq{"id": qps.Exclude})
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return
	}
	roles = []rbac.Item{}
	err = s.db.Select(&roles, sql, args...)
	return
}

func (s *service) Get(id string) (*Role, error) {
	sql, args, err := squirrel.Select("*").From("auth_items").Where(squirrel.Eq{"id": id}).ToSql()
	role := new(Role)
	if err = s.db.Get(role, sql, args...); err != nil {
		return nil, err
	}

	query := squirrel.Select("id").From("auth_items").Where(squirrel.Eq{"item_type": rbac.TypePermission})
	policies, err := s.enforcer.GetImplicitPermissionsForUser(id)
	if err != nil {
		return nil, err
	}
	or := squirrel.Or{}
	for _, v := range policies {
		or = append(or, squirrel.Eq{
			"obj": v[1],
			"act": v[2],
		})
	}
	query = query.Where(or)
	sql, args, err = query.ToSql()
	rows, err := s.db.Queryx(sql, args...)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		role.Permissions = append(role.Permissions, id)
	}

	role.Roles, err = s.enforcer.GetRolesForUser(id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *service) Create(form *Form) (*Role, error) {
	now := time.Now()
	query, args, err := squirrel.Insert("auth_items").SetMap(clevergo.Map{
		"id":         form.ID,
		"item_type":  rbac.TypeRole,
		"name":       form.Name,
		"created_at": now,
		"updated_at": sqlex.ToNullTime(now),
	}).ToSql()
	if err != nil {
		return nil, err
	}

	err = s.db.Transact(func(tx *sql.Tx) error {
		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
		for _, role := range form.Roles {
			if _, err := s.enforcer.AddRoleForUser(form.ID, role); err != nil {
				return err
			}
		}
		permission := new(rbac.Item)
		for _, permissionName := range form.Permissions {
			permissionQuery, args, err := squirrel.Select("*").From("auth_items").Where(squirrel.Eq{
				"id":        permissionName,
				"item_type": rbac.TypePermission,
			}).ToSql()
			if err = s.db.Get(permission, permissionQuery, args...); err != nil {
				return err
			}
			if _, err := s.enforcer.AddPermissionForUser(form.ID, permission.Obj, permission.Act); err != nil {
				return err
			}
		}
		return nil
	})
	return nil, err
}

func (s *service) Delete(id string) error {
	role := new(rbac.Item)
	query, args, err := squirrel.Select("*").From("auth_items").Where(squirrel.Eq{
		"id":        id,
		"item_type": rbac.TypeRole,
	}).ToSql()
	if err := s.db.Get(role, query, args...); err != nil {
		return err
	}
	if role.Reserved {
		return fmt.Errorf("role %s is reserved", role.ID)
	}
	if users, _ := s.enforcer.GetUsersForRole(role.ID); len(users) > 0 {
		return fmt.Errorf("role %s has been assigned to users", role.ID)
	}

	query, args, err = squirrel.Delete("auth_items").Where(squirrel.Eq{
		"id":        id,
		"item_type": rbac.TypeRole,
	}).ToSql()
	if err != nil {
		return err
	}

	err = s.db.Transact(func(tx *sql.Tx) (err error) {
		if _, err = tx.Exec(query, args...); err != nil {
			return
		}
		if _, err = s.enforcer.DeleteRole(id); err != nil {
			return
		}

		return
	})

	return err
}
