package controllers

import (
	"fmt"
	"net/http"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/bootstrap"
	"github.com/clevergo/jsend"
)

func RegisterSite(router clevergo.IRouter, app *frontend.Application, captchaManager *captchas.Manager) {
	s := &site{app, captchaManager}

	router.Get("/", s.index, clevergo.RouteName("index"))
	router.Get("/robots.txt", s.robots)
	router.Get("/about", s.about, clevergo.RouteName("about"))
	router.Get("/contact", s.contact, clevergo.RouteName("contact"))
	router.Post("/contact", s.contact)
}

type site struct {
	*frontend.Application
	captchaManager *captchas.Manager
}

func (s *site) index(ctx *clevergo.Context) error {
	return ctx.Render(http.StatusOK, "site/index.tmpl", nil)
}

func (s *site) about(ctx *clevergo.Context) error {
	return ctx.Render(http.StatusOK, "site/about.tmpl", nil)
}

func (s *site) contact(ctx *clevergo.Context) error {
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

func (s *site) robots(ctx *clevergo.Context) error {
	ctx.WriteString(fmt.Sprintf("User-agent: %s\n", "*"))
	return nil
}
