package web

import (
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/asset"
	"github.com/clevergo/demo/pkg/params"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/clevergo/views/v2"
	"github.com/eko/gocache/store"
	"github.com/go-mail/mail"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	cache          store.StoreInterface
	db             *sqlx.DB
	logger         log.Logger
	sessionManager *scs.SessionManager
	mailer         *mail.Dialer
	userManager    *users.Manager
	captchaManager *captchas.Manager
	params         params.Params
	viewManager    *views.Manager
	beforeRender   func(app *Application, ctx *clevergo.Context, view string, layout bool, viewCtx views.Context)
	accessManager  *access.AccessManager
	assetManager   *asset.AssetManager
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

func (app *Application) ViewManager() *views.Manager {
	return app.viewManager
}

func (app *Application) Logger() log.Logger {
	return app.logger
}

func (app *Application) Mailer() *mail.Dialer {
	return app.mailer
}

func (app *Application) CaptcpaManager() *captchas.Manager {
	return app.captchaManager
}

func (app *Application) SessionManager() *scs.SessionManager {
	return app.sessionManager
}

func (app *Application) UserManager() *users.Manager {
	return app.userManager
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

func (app *Application) Render(ctx *clevergo.Context, view string, viewCtx views.Context) error {
	return app.render(ctx, view, true, viewCtx)
}

func (app *Application) RenderPartial(ctx *clevergo.Context, view string, viewCtx views.Context) error {
	return app.render(ctx, view, false, viewCtx)
}

func (app *Application) render(ctx *clevergo.Context, view string, layout bool, viewCtx views.Context) error {
	ctx.SetContentTypeHTML()
	if viewCtx == nil && app.beforeRender != nil {
		viewCtx = views.Context{}
	}

	if app.beforeRender != nil {
		app.beforeRender(app, ctx, view, layout, viewCtx)
	}

	if layout {
		return app.ViewManager().Render(ctx.Response, view, viewCtx)
	}

	return app.ViewManager().RenderPartial(ctx.Response, view, viewCtx)
}

func (app *Application) Params() params.Params {
	return app.params
}
