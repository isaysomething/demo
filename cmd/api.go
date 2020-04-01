package cmd

import (
	stdlog "log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/api/controllers"
	commonctlrs "github.com/clevergo/demo/internal/controllers"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
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
	provideAPIApp, provideAPIRouteGroups, provideAPIUserManager,
	controllers.Set,
)

func provideAPIServer(
	logger log.Logger,
	routeGroups apiRouteGroups,
	userManager *apiUserManager,
	authenticator auth.Authenticator,
	corsMidware core.CORSMiddleware,
) *core.Server {
	router := clevergo.NewRouter()
	router.Use(
		clevergo.MiddlewareFunc(corsMidware),
	)
	srv := core.NewServer(router, logger)
	srv.Addr = ":4040"
	for _, g := range routeGroups {
		g.Register(router)
	}

	return srv
}

func provideAPIApp(
	logger log.Logger,
	db *sqlx.DB,
	sessionManager *scs.SessionManager,
	userManager *apiUserManager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *api.Application {
	app := newApp(logger, cfg.Params, db, nil, sessionManager, userManager.Manager, mailer, captchaManager, accessManager)
	return &api.Application{Application: app}
}

type apiRouteGroups routeutil.Groups

func provideAPIRouteGroups(
	accessManager *access.Manager,
	post *controllers.Post,
	user *controllers.User,
	captcha *commonctlrs.Captcha,
) apiRouteGroups {
	return apiRouteGroups{
		routeutil.NewGroup("/v1", routeutil.Routes{
			routeutil.NewRoute(http.MethodPost, "/captcha", captcha.Generate),
			routeutil.NewRoute(http.MethodPost, "/check-captcha", captcha.Verify),

			routeutil.NewRoute(http.MethodPost, "/user/login", user.Login),
			routeutil.NewRoute(http.MethodGet, "/user/info", user.Info),
			routeutil.NewRoute(http.MethodPost, "/user/logout", user.Logout),

			routeutil.NewRoute(http.MethodGet, "/users", user.Index).Name("users"),
			routeutil.NewRoute(http.MethodGet, "/users/:id", user.Index).Name("user"),

			routeutil.NewRoute(http.MethodGet, "/posts", post.Index).Name("posts"),
			routeutil.NewRoute(http.MethodGet, "/posts/:id", post.View),
			routeutil.NewRoute(http.MethodPost, "/posts/:id", post.Create),
			routeutil.NewRoute(http.MethodPut, "/posts/:id", post.Update),
			routeutil.NewRoute(http.MethodDelete, "/posts/:id", post.Delete),
		}),
	}
}

type apiUserManager struct {
	*users.Manager
}

func provideAPIUserManager(identityStore auth.IdentityStore) *apiUserManager {
	return &apiUserManager{Manager: users.New(identityStore)}
}
