package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
)

// Post controller.
type Post struct {
	app *api.Application
}

// NewPost returns an Post controller.
func NewPost(app *api.Application) *Post {
	return &Post{app: app}
}

// Index returns posts list.
func (u *Post) Index(ctx *clevergo.Context) error {
	return nil
}

// View returns post info.
func (u *Post) View(ctx *clevergo.Context) error {
	return nil
}

// Create creates a post.
func (u *Post) Create(ctx *clevergo.Context) error {
	return nil
}

// Update udpates a post.
func (u *Post) Update(ctx *clevergo.Context) error {
	return nil
}

// Delete deletes a post.
func (u *Post) Delete(ctx *clevergo.Context) error {
	return nil
}
