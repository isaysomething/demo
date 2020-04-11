package controllers

import (
	"net/http"
	"strconv"

	"github.com/clevergo/form"
	"github.com/clevergo/jsend"

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
func (p *Post) Index(ctx *clevergo.Context) error {
	page := ctx.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	limit := ctx.DefaultQuery("limit", "10")
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		return nil
	}

	count, err := models.GetPostCount(p.DB())
	if err != nil {
		return err
	}

	users, err := models.GetPosts(p.DB(), pageNum, limitNum)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(map[string]interface{}{
		"total": count,
		"items": users,
	}))
}

// View returns post info.
func (p *Post) View(ctx *clevergo.Context) error {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return err
	}
	post, err := models.GetPost(p.DB(), id)
	return ctx.JSON(http.StatusOK, jsend.New(post))
}

// Create creates a post.
func (p *Post) Create(ctx *clevergo.Context) error {
	post := new(models.Post)
	if err := form.Decode(ctx.Request, post); err != nil {
		return err
	}

	if err := post.Save(p.DB()); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(post))
}

// Update udpates a post.
func (p *Post) Update(ctx *clevergo.Context) error {
	return nil
}

// Delete deletes a post.
func (p *Post) Delete(ctx *clevergo.Context) error {
	return nil
}
