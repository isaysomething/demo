package user

import (
	"net/http"
	"strconv"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/jsend"
)

func RegisterRoutes(router clevergo.IRouter, app *api.Application, captchaManager *captchas.Manager, jwtManager *core.JWTManager) {
	r := &resource{app, captchaManager, jwtManager}
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
}

func (r *resource) query(ctx *clevergo.Context) error {
	page := ctx.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	limit := ctx.DefaultQuery("limit", "10")
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	count, err := models.GetUserCount(r.DB())
	if err != nil {
		return err
	}

	users, err := models.GetUsers(r.DB(), pageNum, limitNum)
	if err != nil {
		return err
	}
	return ctx.JSON(200, jsend.New(map[string]interface{}{
		"total": count,
		"items": users,
	}))
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
	identity, _ := user.GetIdentity().(*models.User)
	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"id":    identity.ID,
		"roles": []string{"admin"},
	}))
}

func (r *resource) logout(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	if err := user.Logout(ctx.Request, ctx.Response); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
