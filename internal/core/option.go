package core

import (
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/eko/gocache/store"
	"github.com/go-mail/mail"
	"github.com/jmoiron/sqlx"
)

type Option func(*Application)

func SetDB(db *sqlx.DB) Option {
	return func(app *Application) {
		app.db = db
	}
}

func SetCache(cache store.StoreInterface) Option {
	return func(app *Application) {
		app.cache = cache
	}
}

func SetSessionManager(manager *scs.SessionManager) Option {
	return func(app *Application) {
		app.sessionManager = manager
	}
}

func SetParams(ps Params) Option {
	return func(app *Application) {
		app.params = ps
	}
}

func SetLogger(logger log.Logger) Option {
	return func(app *Application) {
		app.logger = logger
	}
}

func SetUserManager(m *users.Manager) Option {
	return func(app *Application) {
		app.userManager = m
	}
}

func SetAccessManager(manager *access.Manager) Option {
	return func(app *Application) {
		app.accessManager = manager
	}
}

func SetMailer(mailer *mail.Dialer) Option {
	return func(app *Application) {
		app.mailer = mailer
	}
}
