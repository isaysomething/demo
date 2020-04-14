package post

import (
	"net/http"

	"github.com/clevergo/demo/pkg/rest/pagination"

	"github.com/clevergo/jsend"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/models"
)

func RegisterRoutes(router clevergo.IRouter, app *api.Application) {
	r := &resource{app}
	router.Get("/posts", r.query)
	router.Get("/posts/:id", r.get)
	router.Post("/posts", r.create)
	router.Put("/posts/:id", r.update)
	router.Delete("/posts/:id", r.delete)
}

type resource struct {
	*api.Application
}

func (r *resource) query(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)

	p.Items, err = models.GetPosts(r.DB(), p.Limit, p.Offset())
	if err != nil {
		return err
	}
	p.Total, err = models.GetPostsCount(r.DB())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(p))
}

func (r *resource) get(ctx *clevergo.Context) error {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return err
	}
	post, err := models.GetPost(r.DB(), id)
	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *resource) create(ctx *clevergo.Context) error {
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return err
	}
	post, err := form.Create(r.DB())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *resource) update(ctx *clevergo.Context) error {
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return err
	}
	post, err := r.getPost(ctx)
	if err != nil {
		return err
	}
	if err = form.Update(r.DB(), post); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *resource) delete(ctx *clevergo.Context) error {
	post, err := r.getPost(ctx)
	if err != nil {
		return err
	}
	if err := post.Delete(r.DB()); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(nil))
}

func (r *resource) getPost(ctx *clevergo.Context) (*models.Post, error) {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return nil, err
	}
	return models.GetPost(r.DB(), id)
}
