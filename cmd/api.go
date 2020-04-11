package cmd

import (
	stdlog "log"

	"github.com/clevergo/auth"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/api/controllers"
	commonctlrs "github.com/clevergo/demo/internal/controllers"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/clevergo/log"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

var serveAPICmd = &cobra.Command{
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
	provideAPIRouteGroups,
	api.NewUserManager,
	controllers.Set,
)

func provideAPIServer(
	logger log.Logger,
	routeGroups *apiRouteGroups,
	userManager *api.UserManager,
	authenticator auth.Authenticator,
	corsMidware core.CORSMiddleware,
) *core.Server {
	router := clevergo.NewRouter()
	router.Use(
		clevergo.MiddlewareFunc(corsMidware),
		userManager.Middleware(authenticator),
	)
	router.ErrorHandler = api.NewErrorHandler()

	srv := core.NewServer(router, logger)
	srv.Addr = cfg.API.Addr
	routeGroups.Register(router)

	return srv
}

type apiRouteGroups struct {
	routeutil.Groups
}

func provideAPIRouteGroups(
	accessManager *access.Manager,
	post *controllers.Post,
	user *controllers.User,
	captcha *commonctlrs.Captcha,
) *apiRouteGroups {
	gs := routeutil.Groups{
		routeutil.NewGroup("/v1", routeutil.Routes{
			routeutil.Post("/captcha", captcha.Generate),
			routeutil.Post("/check-captcha", captcha.Verify),

			routeutil.Post("/user/login", user.Login),
			routeutil.Get("/user/info", user.Info),
			routeutil.Post("/user/logout", user.Logout),

			routeutil.Get("/users", user.Index).Name("users"),
			routeutil.Get("/users/:id", user.Index).Name("user"),

			routeutil.Get("/posts", post.Index).Name("posts"),
			routeutil.Get("/posts/:id", post.View),
			routeutil.Post("/posts/:id", post.Create),
			routeutil.Put("/posts/:id", post.Update),
			routeutil.Delete("/posts/:id", post.Delete),
		}),
	}

	return &apiRouteGroups{gs}
}
