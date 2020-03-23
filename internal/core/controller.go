package core

type Controller struct {
	*Application
}

func NewController(app *Application) *Controller {
	return &Controller{
		Application: app,
	}
}
