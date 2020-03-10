package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
)

// Site controller.
type Site struct {
	app *backend.Application
}

// NewSite returns a site controller.
func NewSite(app *backend.Application) *Site {
	return &Site{app}
}

// Index displays dashboard.
func (s *Site) Index(ctx *clevergo.Context) error {
	return s.app.Render(ctx, "site/index", nil)
}
