package captcha

import (
	"net/http"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/jsend"
)

type form struct {
	ID      string `json:"id"`
	Captcha string `json:"captcha"`
}

type Captcha struct {
	manager *captchas.Manager
}

func New(manager *captchas.Manager) *Captcha {
	return &Captcha{
		manager: manager,
	}
}

func (c *Captcha) RegisterRoutes(router clevergo.IRouter) {
	router.Post("/captcha", c.generate, clevergo.RouteName("captcha"))
	router.Post("/check-captcha", c.verify)
}

func (c *Captcha) generate(ctx *clevergo.Context) error {
	captcha, err := c.manager.Generate()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"id":   captcha.ID(),
		"data": captcha.EncodeToString(),
	}))
}

func (c *Captcha) verify(ctx *clevergo.Context) error {
	f := form{}
	if err := ctx.Decode(&f); err != nil {
		return err
	}

	if err := c.manager.Verify(f.ID, f.Captcha, false); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
