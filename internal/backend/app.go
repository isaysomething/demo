package backend

import (
	"github.com/clevergo/demo/internal/web"
)

// Application wraps web.Application.
type Application struct {
	*web.Application
}

func New(app *web.Application) *Application {
	return &Application{
		Application: app,
	}
}
