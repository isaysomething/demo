package controllers

import (
	"net/http"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/models"
)

// User controller.
type User struct {
	*api.Application
	captchaManager *captchas.Manager
	jwtManager     *core.JWTManager
}

// NewUser returns an user controller.
func NewUser(app *api.Application, captchaManager *captchas.Manager, jwtManager *core.JWTManager) *User {
	return &User{app, captchaManager, jwtManager}
}

// Index returns users list.
func (u *User) Index(ctx *clevergo.Context) error {
	return nil
}

// View returns user info.
func (u *User) View(ctx *clevergo.Context) error {
	return nil
}

// Login handles login request.
func (u *User) Login(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/backend/", http.StatusFound)
		return nil
	}
	form := forms.NewLogin(u.DB(), user, u.captchaManager)
	v, err := form.Handle(ctx)
	if err != nil {
		return u.Error(ctx, err)
	}

	token, err := u.jwtManager.New(v.GetID())
	if err != nil {
		return u.Error(ctx, err)
	}

	return u.Success(ctx, map[string]interface{}{
		"access_token": token,
	})
}

// Info returns current user info.
func (u *User) Info(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)

	identity, _ := user.GetIdentity().(*models.User)
	return u.Success(ctx, map[string]interface{}{
		"id":    identity.ID,
		"roles": []string{"admin"},
	})
}

// Logout handles logout request.
func (u *User) Logout(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)

	if err := user.Logout(ctx.Request, ctx.Response); err != nil {
		return u.Error(ctx, err)
	}

	return u.Success(ctx, nil)
}
