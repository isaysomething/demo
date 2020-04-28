package authz

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/rbac"
	"github.com/clevergo/demo/pkg/rest/pagination"
	"github.com/clevergo/jsend"
)

type Resource struct {
	*api.Application
	enforcer *casbin.Enforcer
	service  Service
}

func New(app *api.Application, enforcer *casbin.Enforcer, service Service) *Resource {
	return &Resource{
		Application: app,
		enforcer:    enforcer,
		service:     service,
	}
}

func (r *Resource) RegisterRoutes(router clevergo.IRouter) {
	router.Get("/roles", r.queryRoles)
	router.Post("/roles", r.createRole)
	router.Get("/roles/:name", r.getRole)
	router.Put("/roles/:name", r.updateRole)
	router.Delete("/roles/:name", r.deleteRole)
	router.Get("/permissions", r.queryPermissions)
	router.Get("/permission-groups", r.permissionGroups)
}

func (r *Resource) queryRoles(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)
	p.Total, err = r.service.Count()
	if err != nil {
		return err
	}
	p.Items, err = r.service.Query(p.Limit, p.Offset())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(p))
}

func (r *Resource) getRole(ctx *clevergo.Context) error {
	name := ctx.Params.String("name")
	if !r.hasRole(name) {
		ctx.NotFound()
		return nil
	}

	query := squirrel.Select("id").From("auth_items").Where(squirrel.Eq{"item_type": rbac.TypePermission})
	policies, err := r.enforcer.GetImplicitPermissionsForUser(name)
	if err != nil {
		return err
	}
	or := squirrel.Or{}
	for _, v := range policies {
		or = append(or, squirrel.Eq{
			"obj": v[1],
			"act": v[2],
		})
	}
	query = query.Where(or)
	sql, args, err := query.ToSql()
	permissions := []string{}
	rows, err := r.DB().Queryx(sql, args...)
	id := ""
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return err
		}
		permissions = append(permissions, id)
	}

	return ctx.JSON(http.StatusOK, jsend.New(clevergo.Map{
		"name":        name,
		"permissions": permissions,
	}))
}

func (r *Resource) hasRole(name string) bool {
	for _, role := range r.enforcer.GetAllRoles() {
		if role == name {
			return true
		}
	}

	return false
}

func (r *Resource) createRole(ctx *clevergo.Context) error {
	return nil
}

func (r *Resource) updateRole(ctx *clevergo.Context) error {
	return nil
}

var reservedRoles = []string{"admin", "user"}

func isReservedRole(name string) bool {
	for _, role := range reservedRoles {
		if name == role {
			return true
		}
	}

	return false
}

func (r *Resource) deleteRole(ctx *clevergo.Context) error {
	name := ctx.Params.String("name")
	if !r.hasRole(name) {
		ctx.NotFound()
		return nil
	}
	if isReservedRole(name) {
		return fmt.Errorf("role %q is reserved, you cannot delete it", name)
	}

	ok, err := r.enforcer.DeleteRole(name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("failed to delete role %q", name)
	}

	return nil
}

type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var permissions = []Permission{
	{"create_post", ""},
}

func (r *Resource) queryPermissions(ctx *clevergo.Context) error {
	return ctx.JSON(http.StatusOK, jsend.New(permissions))
}

func (r *Resource) permissionGroups(ctx *clevergo.Context) error {
	sql, args, err := squirrel.Select("g.*").From("auth_item_groups g").
		Where("EXISTS (SELECT 1 FROM auth_items WHERE group_id = g.id)").
		ToSql()
	if err != nil {
		return err
	}
	groups := []rbac.Group{}
	if err = r.DB().Select(&groups, sql, args...); err != nil {
		return err
	}
	for i, group := range groups {
		sql, args, err := squirrel.Select("id, name").From("auth_items").
			Where(squirrel.Eq{
				"group_id":  group.ID,
				"item_type": rbac.TypePermission,
			}).
			ToSql()
		if err != nil {
			return err
		}
		if err = r.DB().Select(&groups[i].Permissions, sql, args...); err != nil {
			return err
		}
	}
	return ctx.JSON(http.StatusOK, jsend.New(groups))
}
