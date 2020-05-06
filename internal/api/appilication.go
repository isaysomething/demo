package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/jsend"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Application API application.
type Application struct {
	*core.Application
}

// NewApplication returns API application.
func NewApplication(app *core.Application) *Application {
	return &Application{Application: app}
}

func (app *Application) Success(ctx *clevergo.Context, data interface{}) error {
	return jsend.Success(ctx.Response, data)
}

func (app *Application) Error(ctx *clevergo.Context, err error) error {
	return jsend.Error(ctx.Response, err.Error())
}

func New(
	logger log.Logger,
	params core.Params,
	db *sqlex.DB,
	sessionManager *scs.SessionManager,
	userManager *UserManager,
	mailer *mail.Dialer,
	accessManager *access.Manager,
) *Application {
	boil.SetDB(db)
	boil.DebugMode = true
	opts := []core.Option{
		core.SetParams(params),
		core.SetLogger(logger),
		core.SetDB(db),
		core.SetSessionManager(sessionManager),
		core.SetMailer(mailer),
		core.SetUserManager(userManager.Manager),
		core.SetAccessManager(accessManager),
	}
	return &Application{core.New(opts...)}
}

type UserManager struct {
	*users.Manager
}

func NewUserManager(identityStore auth.IdentityStore) *UserManager {
	return &UserManager{Manager: users.New(identityStore)}
}
