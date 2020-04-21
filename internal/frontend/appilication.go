package frontend

import (
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
)

type Application struct {
	*core.Application
}

func New(
	logger log.Logger,
	params core.Params,
	db *sqlex.DB,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	accessManager *access.Manager,
) *Application {
	opts := []core.Option{
		core.SetParams(params),
		core.SetLogger(logger),
		core.SetDB(db),
		core.SetSessionManager(sessionManager),
		core.SetMailer(mailer),
		core.SetUserManager(userManager),
		core.SetAccessManager(accessManager),
	}
	return &Application{core.New(opts...)}
}
