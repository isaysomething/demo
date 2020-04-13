package controllers

import (
	"net/http"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/form"
	"github.com/clevergo/jsend"
)

func RegisterCaptcha(router clevergo.IRouter, manager *captchas.Manager) {
	c := &captcha{manager}
	router.Post("/captcha", c.generate, clevergo.RouteName("captcha"))
	router.Post("/check-captcha", c.verify)
}

type captchaForm struct {
	ID      string `json:"id"`
	Captcha string `json:"captcha"`
}

type captcha struct {
	manager *captchas.Manager
}

func (c *captcha) generate(ctx *clevergo.Context) error {
	captcha, err := c.manager.Generate()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"id":   captcha.ID(),
		"data": captcha.EncodeToString(),
	}))
}

func (c *captcha) verify(ctx *clevergo.Context) error {
	f := captchaForm{}
	if err := form.Decode(ctx.Request, &f); err != nil {
		return err
	}

	if err := c.manager.Verify(f.ID, f.Captcha, false); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
