package user

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/clevergo/demo/pkg/rest/pagination"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/jsend"
)

func RegisterRoutes(router clevergo.IRouter, app *api.Application, captchaManager *captchas.Manager, jwtManager *core.JWTManager, enforcer *casbin.Enforcer) {
	r := &resource{app, captchaManager, jwtManager, enforcer}
	router.Post("/user/login", r.login)
	router.Get("/user/info", r.info)
	router.Post("/user/logout", r.logout)
	router.Get("/users", r.query)
	router.Get("/users/:id", r.get)
}

type resource struct {
	*api.Application
	captchaManager *captchas.Manager
	jwtManager     *core.JWTManager
	enforcer       *casbin.Enforcer
}

func (r *resource) query(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)
	p.Items, err = models.GetUsers(r.DB(), p.Limit, p.Offset())
	if err != nil {
		return err
	}
	p.Total, err = models.GetUsersCount(r.DB())
	if err != nil {
		return err
	}
	return ctx.JSON(200, jsend.New(p))
}

func (r *resource) get(ctx *clevergo.Context) error {
	return nil
}

func (r *resource) login(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/backend/", http.StatusFound)
		return nil
	}
	form := forms.NewLogin(r.DB(), user, r.captchaManager)
	v, err := form.Handle(ctx)
	if err != nil {
		return err
	}
	token, err := r.jwtManager.New(v.GetID())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"access_token": token,
	}))
}

func (r *resource) info(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	roles, err := r.enforcer.GetRolesForUser("user_" + user.GetIdentity().GetID())
	if err != nil {
		return err
	}
	identity, _ := user.GetIdentity().(*models.User)
	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"id":    identity.ID,
		"roles": roles,
	}))
}

func (r *resource) logout(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	if err := user.Logout(ctx.Request, ctx.Response); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
