package controllers

import (
	"fmt"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/bootstrap"
	"github.com/clevergo/jsend"
)

type Site struct {
	*frontend.Application
}

func NewSite(app *frontend.Application) *Site {
	return &Site{app}
}

func (s *Site) Index(ctx *clevergo.Context) error {
	return s.Render(ctx, "site/index", nil)
}

func (s *Site) About(ctx *clevergo.Context) error {
	return s.Render(ctx, "site/about", nil)
}

func (s *Site) Contact(ctx *clevergo.Context) error {
	captcha, err := s.CaptcpaManager().Generate()
	if err != nil {
		return err
	}

	form := forms.NewContact(s.Mailer(), s.CaptcpaManager())
	if ctx.IsPost() {
		if err = form.Handle(ctx); err == nil {
			s.AddFlash(ctx, bootstrap.NewSuccessAlert("Thanks for contacting us, we'll get in touch with you as soon as possible."))
			ctx.Redirect("/contact", http.StatusFound)
			return nil
		}
		s.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
		s.Logger().Error(err)
	}
	return s.Render(ctx, "site/contact", core.ViewData{
		"form":    form,
		"error":   err,
		"captcha": captcha,
	})
}

func (s *Site) Captcha(ctx *clevergo.Context) error {
	captcha, err := s.CaptcpaManager().Generate()
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	data := map[string]string{
		"id":   captcha.ID(),
		"data": captcha.EncodeToString(),
	}

	return jsend.Success(ctx.Response, data)
}

func (s *Site) CheckCaptcha(ctx *clevergo.Context) error {
	id := ctx.Request.PostFormValue("id")
	captcha := ctx.Request.PostFormValue("captcha")
	err := s.CaptcpaManager().Verify(id, captcha, false)
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}

func (s *Site) Robots(ctx *clevergo.Context) error {
	ctx.WriteString(fmt.Sprintf("User-agent: %s\n", "*"))
	return nil
}
