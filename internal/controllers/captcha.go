package controllers

import (
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/form"
	"github.com/clevergo/jsend"
)

type captcha struct {
	ID string `json:"id"`
	Captcha string `json:"captcha"`
}

type Captcha struct {
	manager *captchas.Manager
}

func NewCaptcha(manager *captchas.Manager) *Captcha {
	return &Captcha{manager: manager}
}

func (c *Captcha) Generate(ctx *clevergo.Context) error {
	captcha, err := c.manager.Generate()
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	data := map[string]string{
		"id":   captcha.ID(),
		"data": captcha.EncodeToString(),
	}

	return jsend.Success(ctx.Response, data)
}

func (c *Captcha) Verify(ctx *clevergo.Context) error {
	captcha := captcha{}
	if err := form.Decode(ctx.Request, &captcha); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}
	
	if err := c.manager.Verify(captcha.ID, captcha.Captcha, false); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}