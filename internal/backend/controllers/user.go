package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
)

// User controller.
type User struct {
	app *backend.Application
}

// NewUser returns a user controller.
func NewUser(app *backend.Application) *User {
	return &User{app}
}

// Index displays users page.
func (u *User) Index(ctx *clevergo.Context) error {
	_, err := u.app.AccessManager().Add(ctx.Params.String("user"), ctx.Params.String("role"))
	if err != nil {
		return err
	}
	return u.app.Render(ctx, "user/index", nil)
}
