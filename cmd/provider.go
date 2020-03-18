package cmd

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"time"

	// database drivers.
	"github.com/CloudyKit/jet/v3"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"

	redissessionstore "github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/auth/authenticators"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/middlewares"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/i18n"
	"github.com/clevergo/log"
	"github.com/clevergo/log/zapadapter"
	"github.com/clevergo/middleware"
	"github.com/go-mail/mail"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/gorilla/csrf"
	sqlxadapter "github.com/memwey/casbin-sqlx-adapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/clevergo/captchas/stores/redisstore"
	"github.com/go-redis/redis/v7"
)

func provideServer(router *clevergo.Router, logger log.Logger, middlewares []func(http.Handler) http.Handler) *web.Server {
	srv := web.NewServer(router, logger)
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

func provideCaptchaManager() *captchas.Manager {
	// redis client.
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
	store := redisstore.New(
		client,
		redisstore.Expiration(10*time.Minute), // captcha expiration, optional.
		redisstore.Prefix("captchas"),         // redis key prefix, optional.
	)
	return web.NewCaptchaManager(store, cfg.Captcha)
}

func provideMailer() *mail.Dialer {
	mailer := mail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password)
	return mailer
}

func provideRouter(
	app *frontend.Application,
	frontendRoutes frontendRoutes,
) *clevergo.Router {
	router := clevergo.NewRouter()
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	router.NotFound = m.Middleware(http.FileServer(packr.New("public", cfg.Server.Root)))

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
	view *web.ViewManager,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *frontend.Application {
	app := newApp(logger, db, view, sessionManager, userManager, mailer, captchaManager, accessManager)
	return &frontend.Application{Application: app}
}

func newApp(
	logger log.Logger,
	db *sqlx.DB,
	view *web.ViewManager,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *web.Application {
	opts := []web.Option{
		web.Params(cfg.Params),
		web.Logger(logger),
		web.DB(db),
		web.SessionManager(sessionManager),
		web.Mailer(mailer),
		web.CaptchaManager(captchaManager),
		web.BeforeRender(func(event *web.BeforeRenderEvent) {
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
		web.UserManager(userManager),
		web.AccessManager(accessManager),
	}
	if view != nil {
		opts = append(opts, web.SetViewManager(view))
	}
	app := web.New(opts...)

	return app
}

func provideMiddlewares(sessionManager *scs.SessionManager, translators *i18n.Translators, userManager *users.Manager, authenticator auth.Authenticator) (v []func(http.Handler) http.Handler, err error) {
	// v = append(v, handlers.RecoveryHandler())
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	if cfg.Server.Gzip {
		v = append(v, middleware.Compress(cfg.Server.GzipLevel))
	}
	v = append(v, m.Middleware)
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
	login := middlewares.LoginCheckerMiddleware(func(r *http.Request, w http.ResponseWriter) bool {
		user, _ := userManager.Get(r, w)
		return user.IsGuest()
	}, middlewares.NewPathSkipper("/*"))
	v = append(v, userManager.Middleware(authenticator), provideI18NMiddleware(translators), login)

	v = append(v, csrf.Protect([]byte("123456"), csrf.Secure(false)))

	return
}

func provideLogger() (log.Logger, func(), error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	sugar, err := config.Build()
	if err != nil {
		return nil, nil, err
	}

	undo := zap.RedirectStdLog(sugar)
	sugar = sugar.WithOptions(zap.AddCallerSkip(1))

	return zapadapter.New(sugar.Sugar()), func() {
		if err := sugar.Sync(); err != nil {
		}

		undo()
	}, nil
}

func provideDB() (*sqlx.DB, func(), error) {
	db, err := web.NewDB(cfg.DB)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
		}
	}, nil
}

func provideView() *web.ViewManager {
	return newView(cfg.View)
}

func newView(cfg web.ViewConfig) *web.ViewManager {
	viewPath := path.Join(cfg.Path, "views")
	box := packr.New(viewPath, viewPath)
	view := web.NewViewManager(box, cfg)
	/*
		funcMap := template.FuncMap{
			"safeHTMLAttr": func(s string) template.HTMLAttr {
				return template.HTMLAttr(s)
			},
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
			"now": time.Now,
		}
		for name, f := range funcMap {
		view.AddFunc(name, f)
		}
	*/
	return view
}

func provideI18N() (*i18n.Translators, error) {
	i18nOpts := []i18n.Option{i18n.Fallback(cfg.I18N.Fallback)}
	translators := i18n.New(i18nOpts...)
	i18nStore := web.NewFileStore(cfg.I18N.Path, i18n.JSONFileDecoder{})
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

func provideIdentityStore(db *sqlx.DB) auth.IdentityStore {
	return web.NewIdentityStore(db)
}

func provideUserManager(identityStore auth.IdentityStore, sessionManager *scs.SessionManager) *users.Manager {
	m := users.New(identityStore)
	m.SetSessionManager(sessionManager)

	return m
}

func provideAuthenticator(identityStore auth.IdentityStore) auth.Authenticator {
	return authenticators.NewComposite(
		authenticators.NewBearerToken("api", identityStore),
		authenticators.NewQueryToken("access_token", identityStore),
	)
}

func provideErrorHandler(app *web.Application) clevergo.ErrorHandler {
	return web.NewErrorHandler(app)
}

func provideSessionManager(store scs.Store) *scs.SessionManager {
	m := web.NewSessionManager(cfg.Session)
	m.Store = store
	return m
}

func provideSessionStore() scs.Store {
	address := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	opts := []redigo.DialOption{
		redigo.DialDatabase(cfg.Redis.Database),
	}
	if cfg.Redis.Password != "" {
		opts = append(opts, redigo.DialPassword(cfg.Redis.Password))
	}
	return redissessionstore.New(redigo.NewPool(func() (redigo.Conn, error) {
		return redigo.Dial("tcp", address, opts...)
	}, 1000))
}
