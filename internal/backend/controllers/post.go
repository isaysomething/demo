package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
)

// Post controller.
type Post struct {
	app *backend.Application
}

// NewPost returns a post controller.
func NewPost(app *backend.Application) *Post {
	return &Post{app}
}

// Index displays posts page.
func (p *Post) Index(ctx *clevergo.Context) error {
	return p.app.Render(ctx, "post/index", nil)
}

// Index displays posts page.
func (p *Post) Create(ctx *clevergo.Context) error {
	return p.app.Render(ctx, "post/create", nil)
}

// Index displays posts page.
func (p *Post) Edit(ctx *clevergo.Context) error {
	return p.app.Render(ctx, "post/edit", nil)
}
