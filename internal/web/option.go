package web

import (
	"strings"

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

type Option func(*Application)

func DB(db *sqlx.DB) Option {
	return func(app *Application) {
		app.db = db
	}
}

func Cache(cache store.StoreInterface) Option {
	return func(app *Application) {
		app.cache = cache
	}
}

func SessionManager(manager *scs.SessionManager) Option {
	return func(app *Application) {
		app.sessionManager = manager
	}
}

func Params(ps params.Params) Option {
	return func(app *Application) {
		app.params = ps
	}
}

func Logger(logger log.Logger) Option {
	return func(app *Application) {
		app.logger = logger
	}
}

func UserManager(m *users.Manager) Option {
	return func(app *Application) {
		app.userManager = m
	}
}

func AccessManager(manager *access.AccessManager) Option {
	return func(app *Application) {
		app.accessManager = manager
	}
}

func AssetManager(manager *asset.AssetManager) Option {
	return func(app *Application) {
		app.assetManager = manager
	}
}

func ViewManager(m *views.Manager) Option {
	return func(app *Application) {
		m.AddFunc("param", func(name string) interface{} {
			val, _ := app.params.Get(name)
			return val
		})
		m.AddFunc("title", strings.Title)
		app.viewManager = m
	}
}

func Mailer(mailer *mail.Dialer) Option {
	return func(app *Application) {
		app.mailer = mailer
	}
}

func CaptchaManager(manager *captchas.Manager) Option {
	return func(app *Application) {
		app.captchaManager = manager
	}
}

func BeforeRender(f func(app *Application, ctx *clevergo.Context, view string, layout bool, data ViewData)) Option {
	return func(app *Application) {
		app.beforeRender = f
	}
}
