package web

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/views/v2"
)

type ErrorHandler struct {
	app      *Application
	viewName string
}

func NewErrorHandler(app *Application) *ErrorHandler {
	return &ErrorHandler{
		app:      app,
		viewName: "site/error",
	}
}

func (eh *ErrorHandler) Handle(ctx *clevergo.Context, err error) {
	eh.app.logger.Errorln(err)

	var errinfo map[string]interface{}
	switch e := err.(type) {
	case clevergo.StatusError:
		errinfo = map[string]interface{}{
			"code":    e.Status(),
			"message": e.Error(),
		}
	default:
		errinfo = map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		}
	}

	if err := eh.app.Render(ctx, eh.viewName, views.Context{
		"error": errinfo,
	}); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}
}

func (eh *ErrorHandler) Error(ctx *clevergo.Context) error {
	return nil
}
