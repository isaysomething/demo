package frontend

import (
	"html/template"
	"reflect"

	"github.com/CloudyKit/jet/v3"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/i18n"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
	"github.com/gorilla/csrf"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	*core.Application
}

func New(
	logger log.Logger,
	params core.Params,
	db *sqlx.DB,
	view *core.ViewManager,
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
		core.SetViewManager(view),
		core.SetBeforeRender(func(event *core.BeforeRenderEvent) {
			user, _ := userManager.Get(event.Context.Request, event.Context.Response)
			event.Data["user"] = user.GetIdentity()

			event.Vars.SetFunc("csrf", func(_ jet.Arguments) reflect.Value {
				return reflect.ValueOf(template.HTML(csrf.TemplateField(event.Context.Request)))
			})

			event.Vars.SetFunc("flashes", func(_ jet.Arguments) reflect.Value {
				return reflect.ValueOf(event.App.Flashes(event.Context))
			})

			translator := i18n.GetTranslator(event.Context.Request)
			event.Vars.SetFunc("T", func(args jet.Arguments) reflect.Value {
				args.RequireNumOfArguments("T", 1, 1)
				return reflect.ValueOf(translator.Sprintf("%m", args.Get(0).String()))
			})
		}),
		core.SetUserManager(userManager),
		core.SetAccessManager(accessManager),
	}
	return &Application{core.New(opts...)}
}
