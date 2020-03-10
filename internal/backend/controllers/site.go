package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
)

type Site struct {
	app *backend.Application
}

func NewSite(app *backend.Application) *Site {
	return &Site{app}
}

func (s *Site) Index(ctx *clevergo.Context) error {
	return s.app.Render(ctx, "site/index", nil)
}
