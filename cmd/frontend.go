package cmd

import (
	"html/template"
	"io"
	"net/http"
	"path"
	"reflect"

	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/form"
	"github.com/clevergo/i18n"
	"github.com/clevergo/jetrenderer"
	"github.com/gorilla/csrf"

	"github.com/clevergo/demo/internal/controllers/captcha"
	"github.com/clevergo/demo/internal/frontend/controllers"

	"github.com/clevergo/captchas"

	"github.com/CloudyKit/jet/v3"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
	"github.com/google/wire"
)

var frontendSet = wire.NewSet(
	frontend.New,
	controllers.NewSite,
	controllers.NewUser,
)

func provideServer(router *clevergo.Router, logger log.Logger) *core.Server {
	srv := core.NewServer(router, logger)
	srv.Addr = cfg.HTTP.Addr
	return srv
}

func provideRouter(
	app *frontend.Application,
	captchaManager *captchas.Manager,
	userManager *users.Manager,
	csrfMidware core.CSRFMiddleware,
	i18nMidware core.I18NMiddleware,
	gzipMidware core.GzipMiddleware,
	sessionMidware core.SessionMiddleware,
	minifyMidware core.MinifyMiddleware,
	loggingMiddleware core.LoggingMiddleware,
	renderer *jetrenderer.Renderer,
	captchaCtl *captcha.Captcha,
	siteCtl *controllers.Site,
	userCtl *controllers.User,
) *clevergo.Router {
	router := clevergo.NewRouter()
	router.NotFound = http.FileServer(packr.New("public", cfg.HTTP.Root))
	router.Renderer = renderer
	router.Decoder = form.New()
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

	renderer.AddGlobalFunc("route", routeFunc)
	renderer.AddGlobalFunc("param", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("param", 1, 1)
		val, _ := app.Params().Get(a.Get(0).String())

		return reflect.ValueOf(val)
	})
	renderer.SetBeforeRender(func(w io.Writer, name string, vars jet.VarMap, data interface{}, ctx *clevergo.Context) error {
		user, _ := userManager.Get(ctx.Request, ctx.Response)
		vars.Set("user", user.GetIdentity())

		vars.SetFunc("csrf", func(_ jet.Arguments) reflect.Value {
			return reflect.ValueOf(template.HTML(csrf.TemplateField(ctx.Request)))
		})

		vars.SetFunc("flashes", func(_ jet.Arguments) reflect.Value {
			return reflect.ValueOf(app.Flashes(ctx))
		})

		translator := i18n.GetTranslator(ctx.Request)
		vars.SetFunc("T", func(args jet.Arguments) reflect.Value {
			args.RequireNumOfArguments("T", 1, 1)
			return reflect.ValueOf(translator.Sprintf("%m", args.Get(0).String()))
		})

		return nil
	})

	router.ServeFiles("/static/*filepath", packr.New("frontend", path.Join(cfg.View.Path, "static")))
	captchaCtl.RegisterRoutes(router)
	siteCtl.RegisterRoutes(router)
	userCtl.RegisterRoutes(router)

	return router
}
