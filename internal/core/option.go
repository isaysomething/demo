package core

import (
	"reflect"
	"time"

	"github.com/CloudyKit/jet/v3"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/captchas"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/params"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
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

func AccessManager(manager *access.Manager) Option {
	return func(app *Application) {
		app.accessManager = manager
	}
}

func SetViewManager(m *ViewManager) Option {
	return func(app *Application) {
		/*m.AddFunc("param", func(name string) interface{} {
			val, _ := app.params.Get(name)
			return val
		})*/
		//m.AddFunc("title", strings.Title)
		m.AddGlobalFunc("param", func(a jet.Arguments) reflect.Value {
			a.RequireNumOfArguments("param", 1, 1)
			val, _ := app.params.Get(a.Get(0).String())

			return reflect.ValueOf(val)
		})
		m.AddGlobalFunc("now", func(_ jet.Arguments) reflect.Value {
			return reflect.ValueOf(time.Now())
		})
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

func BeforeRender(f func(*BeforeRenderEvent)) Option {
	return func(app *Application) {
		app.beforeRender = f
	}
}
