package api

import (
	"fmt"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/jsend"
)

// Application API application.
type Application struct {
	*core.Application
}

// NewApplication returns API application.
func NewApplication(app *core.Application) *Application {
	return &Application{Application: app}
}

func (app *Application) Success(ctx *clevergo.Context, data interface{}) error {
	return jsend.Success(ctx.Response, data)
}

func (app *Application) Error(ctx *clevergo.Context, err error) error {
	fmt.Printf("%+v\n", err)
	return jsend.Error(ctx.Response, err.Error())
}
