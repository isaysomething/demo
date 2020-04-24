package post

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/pkg/rest/pagination"
	"github.com/clevergo/jsend"
)

type Resource struct {
	*api.Application
	service Service
}

func New(app *api.Application, service Service) *Resource {
	return &Resource{
		Application: app,
		service:     service,
	}
}

func (r *Resource) RegisterRoutes(router clevergo.IRouter) {
	router.Get("/posts", r.query)
	router.Get("/posts/:id", r.get)
	router.Post("/posts", r.create)
	router.Put("/posts/:id", r.update)
	router.Delete("/posts/:id", r.delete)
}

func (r *Resource) query(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)

	qps := new(QueryParams)
	if err := api.DecodeQueryParams(qps, ctx); err != nil {
		return err
	}
	p.Items, err = r.service.Query(p.Limit, p.Offset(), qps)
	if err != nil {
		return err
	}
	p.Total, err = r.service.Count()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(p))
}

func (r *Resource) get(ctx *clevergo.Context) error {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return err
	}
	post, err := r.service.Get(id)
	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *Resource) create(ctx *clevergo.Context) error {
	post, err := r.service.Create(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *Resource) update(ctx *clevergo.Context) error {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return err
	}
	form := new(Form)
	if err := ctx.Decode(form); err != nil {
		return err
	}
	post, err := r.service.Update(id, form)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(post))
}

func (r *Resource) delete(ctx *clevergo.Context) error {
	id, err := ctx.Params.Int64("id")
	if err != nil {
		return err
	}
	err = r.service.Delete(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
