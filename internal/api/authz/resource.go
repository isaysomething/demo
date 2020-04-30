package authz

import (
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
	router.Get("/roles/:id", r.getRole)
	router.Put("/roles/:id", r.updateRole)
	router.Delete("/roles/:id", r.deleteRole)
	router.Get("/permissions", r.queryPermissions)
	router.Get("/permission-groups", r.permissionGroups)
}

func (r *Resource) queryRoles(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)
	p.Total, err = r.service.Count()
	if err != nil {
		return err
	}
	qps := new(QueryParams)
	if err = api.DecodeQueryParams(qps, ctx); err != nil {
		return err
	}
	p.Items, err = r.service.Query(p.Limit, p.Offset(), qps)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(p))
}

func (r *Resource) getRole(ctx *clevergo.Context) error {
	id := ctx.Params.String("id")
	role, err := r.service.Get(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(role))
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
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return err
	}
	if err := form.Validate(); err != nil {
		return err
	}
	role, err := r.service.Create(form)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(role))
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
	id := ctx.Params.String("id")
	if err := r.service.Delete(id); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(nil))
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
