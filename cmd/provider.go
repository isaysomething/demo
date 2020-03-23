package cmd

import (
	"html/template"
	"net/http"
	"path"
	"reflect"

	// database drivers.
	"github.com/CloudyKit/jet/v3"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/middlewares"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/i18n"
	"github.com/clevergo/log"
	"github.com/clevergo/middleware"
	"github.com/go-mail/mail"
	"github.com/gorilla/csrf"
	sqlxadapter "github.com/memwey/casbin-sqlx-adapter"
)

func provideServer(router *clevergo.Router, logger log.Logger, middlewares []func(http.Handler) http.Handler) *core.Server {
	srv := core.NewServer(router, logger)
	srv.Addr = cfg.Server.Addr
	srv.Use(middlewares...)
	return srv
}

func provideEnforcer() (*casbin.Enforcer, error) {
	//return casbin.NewEnforcer("casbin/model.conf", "casbin/policy.csv")
	opts := &sqlxadapter.AdapterOptions{
		DriverName:     "mysql",
		DataSourceName: "root:123456@tcp(127.0.0.1:3306)/clevergo",
		TableName:      "auth_rules",
		// or reuse an existing connection:
		// DB: myDBConn,
	}
	conf := packr.New("rbac", "../conf")
	content, err := conf.FindString("rbac_model.conf")
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(content)
	if err != nil {
		return nil, err
	}

	a := sqlxadapter.NewAdapterFromOptions(opts)
	e, err := casbin.NewEnforcer()
	if err != nil {
		return nil, err
	}
	if err = e.InitWithModelAndAdapter(m, a); err != nil {
		return nil, err
	}

	// Reload the policy from file/database.
	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	// Save the current policy (usually after changed with Casbin API) back to file/database.
	//e.SavePolicy()

	return e, nil
}

func provideRouter(
	app *frontend.Application,
	frontendRoutes frontendRoutes,
) *clevergo.Router {
	router := clevergo.NewRouter()
	router.NotFound = http.FileServer(packr.New("public", cfg.Server.Root))

	router.Use(
		clevergo.Recovery(true),
		middlewares.Logging(core.LoggerWriter(app.Logger())),
		middlewares.Minify(),
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

func provideMiddlewares(sessionManager *scs.SessionManager, translators *i18n.Translators, userManager *users.Manager, authenticator auth.Authenticator) (v []func(http.Handler) http.Handler, err error) {
	// v = append(v, handlers.RecoveryHandler())
	if cfg.Server.Gzip {
		v = append(v, middleware.Compress(cfg.Server.GzipLevel))
	}
	/*
		if cfg.Server.AccessLog {
			var accessLog io.Writer = os.Stdout
			if cfg.Server.AccessLogFile != "" {
				accessLog, err = os.OpenFile(cfg.Server.AccessLogFile, os.O_CREATE|os.O_APPEND, os.FileMode(cfg.Server.AccessLogFileMode))
				if err != nil {
					return
				}
			}
			v = append(v, middleware.Logging(accessLog))
		}
	*/
	login := middlewares.LoginCheckerMiddleware(func(r *http.Request, w http.ResponseWriter) bool {
		user, _ := userManager.Get(r, w)
		return user.IsGuest()
	}, middlewares.NewPathSkipper("/*"))
	v = append(v, userManager.Middleware(authenticator), provideI18NMiddleware(translators), login)

	v = append(v, csrf.Protect([]byte("123456"), csrf.Secure(false)))

	return
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

func provideI18N() (*i18n.Translators, error) {
	i18nOpts := []i18n.Option{i18n.Fallback(cfg.I18N.Fallback)}
	translators := i18n.New(i18nOpts...)
	i18nStore := core.NewFileStore(cfg.I18N.Path, i18n.JSONFileDecoder{})
	if err := translators.Import(i18nStore); err != nil {
		return nil, err
	}

	return translators, nil
}

func provideI18NMiddleware(translators *i18n.Translators) func(http.Handler) http.Handler {
	var languageParsers []i18n.LanguageParser
	if cfg.I18N.Param != "" {
		languageParsers = append(languageParsers, i18n.NewURLLanguageParser(cfg.I18N.Param))
	}
	if cfg.I18N.CookieParam != "" {
		languageParsers = append(languageParsers, i18n.NewCookieLanguageParser(cfg.I18N.CookieParam))
	}
	languageParsers = append(languageParsers, i18n.HeaderLanguageParser{})
	return i18n.Middleware(translators, languageParsers...)
}
