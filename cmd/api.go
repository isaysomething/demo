package cmd

import (
	stdlog "log"

	"github.com/clevergo/captchas"
	"github.com/clevergo/demo/internal/api/authz"
	"github.com/clevergo/demo/internal/api/post"
	"github.com/clevergo/demo/internal/api/user"
	"github.com/clevergo/demo/internal/controllers/captcha"
	"github.com/clevergo/form"

	"github.com/clevergo/auth"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/log"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start API service",
	Long:  `Start API service`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, f, err := initializeAPIServer()
		if err != nil {
			stdlog.Fatal(err.Error())
		}
		defer f()

		if err := srv.ListenAndServe(); err != nil {
			stdlog.Fatal(err.Error())
		}
	},
}

var apiSet = wire.NewSet(
	api.New,
	api.NewUserManager,
	api.NewAuthzMiddleware,

	user.New,
	post.Set,
	authz.New,
)

func provideAPIServer(
	logger log.Logger,
	app *api.Application,
	captchaManager *captchas.Manager,
	jwtManager *core.JWTManager,
	userManager *api.UserManager,
	authenticator auth.Authenticator,
	corsMidware core.CORSMiddleware,
	authzMidware api.AuthzMiddleware,
	userResource *user.Resource,
	postResource *post.Resource,
	authzResource *authz.Resource,
	captchaCtl *captcha.Captcha,
) *core.Server {
	router := clevergo.NewRouter()
	router.Decoder = form.New()
	router.Use(
		clevergo.MiddlewareFunc(corsMidware),
		userManager.Middleware(authenticator),
		clevergo.MiddlewareFunc(authzMidware),
	)
	router.ErrorHandler = api.NewErrorHandler()

	v1 := router.Group("/v1")
	postResource.RegisterRoutes(v1)
	userResource.RegisterRoutes(v1)
	authzResource.RegisterRoutes(v1)
	captchaCtl.RegisterRoutes(v1)

	srv := core.NewServer(router, logger)
	srv.Addr = cfg.API.Addr

	return srv
}
