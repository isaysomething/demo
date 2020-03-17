package api

import (
	"fmt"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/jsend"
)

// Application API application.
type Application struct {
	*web.Application
}

// NewApplication returns API application.
func NewApplication(app *web.Application) *Application {
	return &Application{Application: app}
}

func (app *Application) Success(ctx *clevergo.Context, data interface{}) error {
	return jsend.Success(ctx.Response, data)
}

func (app *Application) Error(ctx *clevergo.Context, err error) error {
	fmt.Printf("%+v\n", err)
	return jsend.Error(ctx.Response, err.Error())
}
