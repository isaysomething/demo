package controllers

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
	"github.com/clevergo/jsend"
)

// Site controller.
type Site struct {
	*backend.Application
}

// NewSite returns a site controller.
func NewSite(app *backend.Application) *Site {
	return &Site{Application: app}
}

// Index displays dashboard.
func (s *Site) Index(ctx *clevergo.Context) error {
	return s.Render(ctx, "site/index", nil)
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
