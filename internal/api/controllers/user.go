package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
)

// User controller.
type User struct {
	app *api.Application
}

// NewUser returns an user controller.
func NewUser(app *api.Application) *User {
	return &User{app: app}
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
	return nil
}

// Logout handles logout request.
func (u *User) Logout(ctx *clevergo.Context) error {
	return nil
}
