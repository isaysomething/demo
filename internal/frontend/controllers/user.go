package controllers

import (
	"github.com/clevergo/demo/internal/frontend"
)

type User struct {
	*frontend.Application
}

func NewUser(app *frontend.Application) *User {
	return &User{app}
}
