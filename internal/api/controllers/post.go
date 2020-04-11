package controllers

import (
	"fmt"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/models"
)

// Post controller.
type Post struct {
	*api.Application
}

// NewPost returns an Post controller.
func NewPost(app *api.Application) *Post {
	return &Post{app}
}

// Index lists posts.
func (u *Post) Index(ctx *clevergo.Context) error {
	posts, err := models.GetPosts(u.DB(), 0, 10)
	fmt.Println(posts, err)
	if err != nil {
		return err
	}

	return u.Success(ctx, posts)
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
