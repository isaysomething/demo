package cmd

import (
	"net/http"
	"path"
	"reflect"

	"github.com/clevergo/captchas"
	commonctlrs "github.com/clevergo/demo/internal/controllers"

	"github.com/CloudyKit/jet/v3"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/frontend/controllers"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/clevergo/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
	"github.com/google/wire"
)

var frontendSet = wire.NewSet(
	frontend.New,
	core.NewViewManager,
	provideRoutes,
	controllers.Set,
)

type routes struct {
	routeutil.Routes
}

func provideRoutes(
	site *controllers.Site,
	user *controllers.User,
) *routes {
	rs := routeutil.Routes{
		routeutil.Get("/", site.Index).Name("home"),
		routeutil.Get("/robots.txt", site.Robots),
		routeutil.Get("/about", site.About).Name("about"),
		routeutil.Get("/contact", site.Contact).Name("contact"),
		routeutil.Post("/contact", site.Contact),

		routeutil.Get("/user", user.Index).Name("user"),
		routeutil.Get("/login", user.Login).Name("login"),
		routeutil.Post("/login", user.Login),
		routeutil.Post("/logout", user.Logout).Name("logout"),
		routeutil.Get("/signup", user.Signup).Name("signup"),
		routeutil.Post("/signup", user.Signup),
		routeutil.Post("/user/check-email", user.CheckEmail),
		routeutil.Post("/user/check-username", user.CheckUsername),
		routeutil.Get("/user/request-reset-password", user.RequestResetPassword).Name("request-reset-password"),
		routeutil.Post("/user/request-reset-password", user.RequestResetPassword),
		routeutil.Get("/user/reset-password", user.ResetPassword).Name("reset-password"),
		routeutil.Post("/user/reset-password", user.ResetPassword),
		routeutil.Get("/user/verify-email", user.VerifyEmail).Name("verify-email"),
		routeutil.Get("/user/resend-verification-email", user.ResendVerificationEmail).Name("resend-verification-email"),
		routeutil.Post("/user/resend-verification-email", user.ResendVerificationEmail),
		routeutil.Post("/user/change-password", user.ChangePassword).Name("change-password"),
	}

	return &routes{rs}
}

func provideServer(router *clevergo.Router, logger log.Logger) *core.Server {
	srv := core.NewServer(router, logger)
	srv.Addr = cfg.HTTP.Addr
	return srv
}

func provideRouter(
	app *frontend.Application,
	routes *routes,
	captchaManager *captchas.Manager,
	csrfMidware core.CSRFMiddleware,
	i18nMidware core.I18NMiddleware,
	gzipMidware core.GzipMiddleware,
	sessionMidware core.SessionMiddleware,
	minifyMidware core.MinifyMiddleware,
	loggingMiddleware core.LoggingMiddleware,
) *clevergo.Router {
	router := clevergo.NewRouter()
	router.NotFound = http.FileServer(packr.New("public", cfg.HTTP.Root))

	router.Use(
		clevergo.Recovery(true),
		clevergo.MiddlewareFunc(loggingMiddleware),
		clevergo.MiddlewareFunc(gzipMidware),
		clevergo.MiddlewareFunc(minifyMidware),
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
	commonctlrs.RegisterCaptcha(router, captchaManager)
	routes.Register(router)

	return router
}
