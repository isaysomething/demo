package cmd

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	// database drivers.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	redissessionstore "github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/auth/authenticators"
	"github.com/clevergo/captchas"
	"github.com/clevergo/captchas/drivers"
	"github.com/clevergo/captchas/memstore"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/backend"
	"github.com/clevergo/demo/internal/common"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/asset"
	"github.com/clevergo/demo/pkg/middlewares"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/i18n"
	"github.com/clevergo/log"
	"github.com/clevergo/log/zapadapter"
	"github.com/clevergo/middleware"
	"github.com/clevergo/views/v2"
	"github.com/go-mail/mail"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

func provideServer(router *clevergo.Router, logger log.Logger, middlewares []func(http.Handler) http.Handler) *web.Server {
	srv := web.NewServer(router, logger)
	srv.Addr = cfg.Addr
	srv.Use(middlewares...)
	return srv
}

func provideCaptchaManager() *captchas.Manager {
	manager := captchas.New(
		memstore.New(),
		drivers.NewDigit(),
	)
	return manager
}

func provideMailer() *mail.Dialer {
	mailer := mail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password)
	return mailer
}

func provideRouter(
	app *frontend.Application,
	frontendRoutes frontendRoutes,
	backendApp *backend.Application,
	backendRoutes backendRoutes,
) *clevergo.Router {
	router := clevergo.NewRouter()
	router.NotFound = http.FileServer(http.Dir(cfg.Root))

	urlFunc := func(name string, args ...string) string {
		url, err := router.URL(name, args...)
		if err != nil {
			app.Logger().Infoln(err)
			return ""
		}
		return url.String()
	}

	app.ViewManager().AddFunc("url", urlFunc)
	backendApp.ViewManager().AddFunc("url", urlFunc)

	router.ServeFiles("/static/*filepath", http.Dir(path.Join(cfg.View.Directory, "static")))
	for _, route := range frontendRoutes {
		route.Register(router)
	}

	router.ServeFiles("/console/static/*filepath", http.Dir(path.Join(cfg.AdminView.Directory, "static")))
	console := router.Group("/console")
	for _, route := range backendRoutes {
		route.Register(console)
	}

	return router
}

func provideApp(
	logger log.Logger,
	db *sqlx.DB,
	view *views.Manager,
	sessionManager *scs.SessionManager,
	userStore *users.Store,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
) *frontend.Application {
	app := newApp(logger, db, view, sessionManager, userStore, mailer, captchaManager)
	return &frontend.Application{app}
}

func provideBackendApp(
	logger log.Logger,
	db *sqlx.DB,
	view *AdminView,
	sessionManager *scs.SessionManager,
	userStore *users.Store,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
) *backend.Application {
	app := newApp(logger, db, view.Manager, sessionManager, userStore, mailer, captchaManager)
	return &backend.Application{app}
}

func newApp(
	logger log.Logger,
	db *sqlx.DB,
	view *views.Manager,
	sessionManager *scs.SessionManager,
	userStore *users.Store,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
) *web.Application {
	opts := []web.Option{
		web.Params(cfg.Params),
		web.Logger(logger),
		web.DB(db),
		web.SessionManager(sessionManager),
		web.AssetManager(&asset.AssetManager{}),
		web.Mailer(mailer),
		web.CaptchaManager(captchaManager),
		web.BeforeRender(func(app *web.Application, ctx *clevergo.Context, view string, layout bool, viewCtx views.Context) {
			translator := i18n.GetTranslator(ctx.Request)
			viewCtx["translator"] = translator
			viewCtx["translate"] = func(key string) string {
				return translator.Sprintf("%m", key)
			}
			user, _ := userStore.Get(ctx.Request, ctx.Response)
			viewCtx["user"] = user.GetIdentity()
			viewCtx["flashes"] = app.Flashes(ctx)
			viewCtx["csrf"] = csrf.TemplateField(ctx.Request)
		}),
		web.UserStore(userStore),
	}
	if view != nil {
		opts = append(opts, web.ViewManager(view))
	}
	app := web.New(opts...)

	return app
}

func provideMiddlewares(sessionManager *scs.SessionManager, translators *i18n.Translators, userStore *users.Store) (v []func(http.Handler) http.Handler, err error) {
	v = append(v, handlers.RecoveryHandler())
	if cfg.Gzip {
		v = append(v, middleware.Compress(cfg.GzipLevel))
	}
	if cfg.AccessLog {
		var accessLog io.Writer = os.Stdout
		if cfg.AccessLogFile != "" {
			accessLog, err = os.OpenFile(cfg.AccessLogFile, os.O_CREATE|os.O_APPEND, os.FileMode(cfg.AccessLogFileMode))
			if err != nil {
				return
			}
		}
		v = append(v, middleware.Logging(accessLog))
	}
	login := middlewares.LoginCheckerMiddleware(func(r *http.Request, w http.ResponseWriter) bool {
		user, _ := userStore.Get(r, w)
		return user.IsGuest()
	}, middlewares.NewPathSkipper("/*"))
	v = append(v, sessionManager.LoadAndSave, provideI18NMiddleware(translators), login)

	v = append(v, csrf.Protect([]byte("123456"), csrf.Secure(false)))

	return
}

func provideLogger() (log.Logger, func(), error) {
	sugar, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return nil, nil, err
	}

	undo := zap.RedirectStdLog(sugar)

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

func provideView() *views.Manager {
	return newView(cfg.View)
}

type AdminView struct {
	*views.Manager
}

func provideAdminView(logger log.Logger) *AdminView {
	view := newView(cfg.AdminView)
	return &AdminView{view}
}

func newView(cfg web.ViewConfig) *views.Manager {
	view := web.NewView(cfg)
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
	return view
}

func provideI18N() (*i18n.Translators, error) {
	i18nOpts := []i18n.Option{i18n.Fallback(cfg.I18N.Fallback)}
	translators := i18n.New(i18nOpts...)
	i18nStore := i18n.NewFileStore(cfg.I18N.Directory, i18n.JSONFileDecoder{})
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
	return common.NewIdentityStore(db)
}

func provideUserStore(identityStore auth.IdentityStore, sessionManager *scs.SessionManager) *users.Store {
	userStore := users.NewStore(identityStore)
	userStore.SetSessionManager(sessionManager)

	return userStore
}

func provideAuthenticator(identityStore auth.IdentityStore) auth.Authenticator {
	return authenticators.NewQueryToken("token", identityStore)
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
	address := "localhost:6379"
	return redissessionstore.New(redigo.NewPool(func() (redigo.Conn, error) {
		return redigo.Dial("tcp", address)
	}, 1000))
}
