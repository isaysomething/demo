package controllers

import (
	"fmt"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/bootstrap"
	"github.com/clevergo/jsend"
)

type Site struct {
	*frontend.Application
	captchaManager *captchas.Manager
}

func NewSite(app *frontend.Application, captchaManager *captchas.Manager) *Site {
	return &Site{app, captchaManager}
}

func (s *Site) Index(ctx *clevergo.Context) error {
	return s.Render(ctx, "site/index", nil)
}

func (s *Site) About(ctx *clevergo.Context) error {
	return s.Render(ctx, "site/about", nil)
}

func (s *Site) Contact(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewContact(s.Mailer(), s.captchaManager)
		if err := form.Handle(ctx); err != nil {
			s.Logger().Error(err)
			return jsend.Error(ctx.Response, err.Error())
		}
		s.AddFlash(ctx, bootstrap.NewSuccessAlert("Thanks for contacting us, we'll get in touch with you as soon as possible."))
		return jsend.Success(ctx.Response, nil)
	}
	err := s.Render(ctx, "site/contact", core.ViewData{})
	fmt.Println(err)
	return err
}

func (s *Site) Robots(ctx *clevergo.Context) error {
	ctx.WriteString(fmt.Sprintf("User-agent: %s\n", "*"))
	return nil
}
