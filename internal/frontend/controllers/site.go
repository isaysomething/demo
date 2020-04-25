package controllers

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/bootstrap"
	"github.com/clevergo/demo/pkg/rest/pagination"
	"github.com/clevergo/jsend"
)

type Site struct {
	*frontend.Application
	captchaManager *captchas.Manager
}

func NewSite(app *frontend.Application, captchaManager *captchas.Manager) *Site {
	return &Site{
		Application:    app,
		captchaManager: captchaManager,
	}
}

func (s *Site) RegisterRoutes(router clevergo.IRouter) {
	router.Get("/", s.index, clevergo.RouteName("index"))
	router.Get("/robots.txt", s.robots)
	router.Get("/about", s.about, clevergo.RouteName("about"))
	router.Get("/contact", s.contact, clevergo.RouteName("contact"))
	router.Post("/contact", s.contact)
}

func (s *Site) index(ctx *clevergo.Context) error {
	pagination := pagination.NewFromContext(ctx)
	sql, args, err := squirrel.Select("count(*)").
		Where(squirrel.Eq{"state": models.PostStatePublished}).
		From("posts").
		ToSql()
	if err != nil {
		return err
	}
	if err = s.DB().Get(&pagination.Total, sql, args...); err != nil {
		return err
	}

	sql, args, err = squirrel.Select("*").
		Where(squirrel.Eq{"state": models.PostStatePublished}).
		From("posts").
		OrderBy("created_at DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		ToSql()
	if err != nil {
		return err
	}
	posts := []models.Post{}
	if err = s.DB().Select(&posts, sql, args...); err != nil {
		return err
	}
	pagination.Items = posts
	return ctx.Render(http.StatusOK, "site/index.tmpl", clevergo.Map{
		"pagination": pagination,
	})
}

func (s *Site) about(ctx *clevergo.Context) error {
	return ctx.Render(http.StatusOK, "site/about.tmpl", nil)
}

func (s *Site) contact(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewContact(s.Mailer(), s.captchaManager)
		if err := form.Handle(ctx); err != nil {
			s.Logger().Error(err)
			return jsend.Error(ctx.Response, err.Error())
		}
		s.AddFlash(ctx, bootstrap.NewSuccessAlert("Thanks for contacting us, we'll get in touch with you as soon as possible."))
		return jsend.Success(ctx.Response, nil)
	}
	return ctx.Render(http.StatusOK, "site/contact.tmpl", core.ViewData{})
}

func (s *Site) robots(ctx *clevergo.Context) error {
	ctx.WriteString(fmt.Sprintf("User-agent: %s\n", "*"))
	return nil
}
