package core

import (
	"github.com/CloudyKit/jet/v3"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/eko/gocache/store"
	"github.com/go-mail/mail"
	"github.com/jmoiron/sqlx"
)

type BeforeRenderEvent struct {
	App     *Application
	Context *clevergo.Context
	View    string
	Vars    jet.VarMap
	Data    ViewData
}

type Application struct {
	cache          store.StoreInterface
	db             *sqlx.DB
	logger         log.Logger
	sessionManager *scs.SessionManager
	mailer         *mail.Dialer
	userManager    *users.Manager
	params         Params
	accessManager  *access.Manager
}

func New(opts ...Option) *Application {
	app := &Application{}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *Application) DB() *sqlx.DB {
	return app.db
}

func (app *Application) Cache() store.StoreInterface {
	return app.cache
}

func (app *Application) Logger() log.Logger {
	return app.logger
}

func (app *Application) Mailer() *mail.Dialer {
	return app.mailer
}

func (app *Application) SessionManager() *scs.SessionManager {
	return app.sessionManager
}

func (app *Application) UserManager() *users.Manager {
	return app.userManager
}

func (app *Application) AccessManager() *access.Manager {
	return app.accessManager
}

func (app *Application) User(ctx *clevergo.Context) (*users.User, error) {
	return app.userManager.Get(ctx.Request, ctx.Response)
}

func (app *Application) Flashes(ctx *clevergo.Context) (flashes Flashes) {
	flashes, _ = app.sessionManager.Pop(ctx.Request.Context(), "_flash").(Flashes)
	return
}

func (app *Application) AddFlash(ctx *clevergo.Context, flash Flash) {
	flashes := append(app.Flashes(ctx), flash)
	app.sessionManager.Put(ctx.Request.Context(), "_flash", flashes)
}

func (app *Application) Params() Params {
	return app.params
}
