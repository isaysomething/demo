package rest

import "github.com/clevergo/demo/internal/core"

type Controller struct {
	*core.Application
}

func NewController(app *core.Application) *Controller {
	return &Controller{
		Application: app,
	}
}
