package controllers

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/form"
	"github.com/clevergo/jsend"
)

// User controller.
type User struct {
	*api.Application
}

// NewUser returns an user controller.
func NewUser(app *api.Application) *User {
	return &User{app}
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
	form := forms.NewLogin(u.DB(), user, u.CaptcpaManager())
	v, err := form.Handle(ctx)
	if err != nil {
		return err
	}

	return jsend.Success(ctx.Response, v)
}

func (u *User) CheckUsername(ctx *clevergo.Context) error {
	f := forms.NewCheckUsername(u.DB())
	err := form.Decode(ctx.Request, f)
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}
	if err = f.Validate(); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}

func (u *User) CheckEmail(ctx *clevergo.Context) error {
	f := forms.NewCheckUserEmail(u.DB())
	err := form.Decode(ctx.Request, f)
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}
	if err = f.Validate(); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}

// Logout handles logout request.
func (u *User) Logout(ctx *clevergo.Context) error {
	return nil
}
