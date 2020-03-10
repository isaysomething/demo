package rest

import "github.com/clevergo/demo/internal/web"

type Controller struct {
	*web.Application
}

func NewController(app *web.Application) *Controller {
	return &Controller{
		Application: app,
	}
}
