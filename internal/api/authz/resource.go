package authz

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/jsend"
)

type Resource struct {
	*api.Application
	enforcer *casbin.Enforcer
}

func New(app *api.Application, enforcer *casbin.Enforcer) *Resource {
	return &Resource{
		Application: app,
		enforcer:    enforcer,
	}
}

func (r *Resource) RegisterRoutes(router clevergo.IRouter) {
	router.Get("/roles", r.queryRoles)
	router.Post("/roles", r.createRole)
	router.Put("/roles/:name", r.updateRole)
	router.Delete("/roles/:name", r.deleteRole)
}

func (r *Resource) queryRoles(ctx *clevergo.Context) error {
	roles := r.enforcer.GetAllRoles()
	return ctx.JSON(http.StatusOK, jsend.New(roles))
}

func (r *Resource) createRole(ctx *clevergo.Context) error {
	return nil
}

func (r *Resource) updateRole(ctx *clevergo.Context) error {
	return nil
}

func (r *Resource) deleteRole(ctx *clevergo.Context) error {
	return nil
}
