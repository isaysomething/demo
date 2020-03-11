package api

import "github.com/clevergo/demo/internal/web"

// Application API application.
type Application struct {
	*web.Application
}

// NewApplication returns API application.
func NewApplication(app *web.Application) *Application {
	return &Application{Application: app}
}
