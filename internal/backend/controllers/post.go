package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
)

type Post struct {
	app *backend.Application
}

func NewPost(app *backend.Application) *Post {
	return &Post{app}
}

func (p *Post) Index(ctx *clevergo.Context) error {
	return p.app.Render(ctx, "post/index", nil)
}
