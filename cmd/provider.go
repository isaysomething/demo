package cmd

import (
	"html/template"
	"net/http"
	"path"
	"reflect"

	"github.com/CloudyKit/jet/v3"
	// database drivers.
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	"github.com/tdewolff/minify/v2"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/i18n"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
	"github.com/gorilla/csrf"
)

func provideServer(router *clevergo.Router, logger log.Logger, m *minify.M) *core.Server {
	srv := core.NewServer(router, logger)
	srv.Addr = cfg.Server.Addr
	srv.Use(
		m.Middleware,
	)
	return srv
}

func provideRouter(
	app *frontend.Application,
	frontendRoutes frontendRoutes,
	csrfMidware core.CSRFMiddleware,
	i18nMidware core.I18NMiddleware,
	gzipMidware core.GzipMiddleware,
	sessionMidware core.SessionMiddleware,
	minifyMidware core.MinifyMiddleware,
	loggingMiddleware core.LoggingMiddleware,
) *clevergo.Router {
	router := clevergo.NewRouter()
	router.NotFound = http.FileServer(packr.New("public", cfg.Server.Root))

	router.Use(
		clevergo.Recovery(true),
		clevergo.MiddlewareFunc(loggingMiddleware),
		clevergo.MiddlewareFunc(minifyMidware),
		clevergo.MiddlewareFunc(gzipMidware),
		clevergo.MiddlewareFunc(sessionMidware),
		clevergo.MiddlewareFunc(csrfMidware),
		clevergo.MiddlewareFunc(i18nMidware),
	)

	routeFunc := func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("route", 1, 1)
		var a []string
		for i := 1; i < args.NumOfArguments(); i++ {
			a = append(a, args.Get(i).String())
		}
		url, err := router.URL(args.Get(0).String(), a...)
		if err != nil {
			app.Logger().Infoln(err)
			return reflect.ValueOf("")
		}
		return reflect.ValueOf(url.String())
	}

	app.ViewManager().AddGlobalFunc("route", routeFunc)

	router.ServeFiles("/static/*filepath", packr.New("frontend", path.Join(cfg.View.Path, "static")))
	for _, route := range frontendRoutes {
		route.Register(router)
	}

	return router
}

func provideApp(
	logger log.Logger,
	db *sqlx.DB,
	view *core.ViewManager,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *frontend.Application {
	app := newApp(logger, cfg.Params, db, view, sessionManager, userManager, mailer, captchaManager, accessManager)
	return &frontend.Application{Application: app}
}

func newApp(
	logger log.Logger,
	params core.Params,
	db *sqlx.DB,
	view *core.ViewManager,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *core.Application {
	opts := []core.Option{
		core.SetParams(params),
		core.SetLogger(logger),
		core.SetDB(db),
		core.SetSessionManager(sessionManager),
		core.SetMailer(mailer),
		core.SetCaptchaManager(captchaManager),
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
	if view != nil {
		opts = append(opts, core.SetViewManager(view))
	}
	app := core.New(opts...)

	return app
}

func provideView() *core.ViewManager {
	return newView(cfg.View)
}

func newView(cfg core.ViewConfig) *core.ViewManager {
	viewPath := path.Join(cfg.Path, "views")
	box := packr.New(viewPath, viewPath)
	view := core.NewViewManager(box, cfg)
	return view
}
