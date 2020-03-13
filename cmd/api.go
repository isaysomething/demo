package cmd

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/captchas"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/api/controllers"
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/log"
	"github.com/go-mail/mail"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var apiSet = wire.NewSet(
	provideAPIApp, provideAPIRoutes,
	controllers.NewUser, controllers.NewPost,
)

func provideAPIApp(
	logger log.Logger,
	db *sqlx.DB,
	view *web.ViewManager,
	sessionManager *scs.SessionManager,
	userManager *users.Manager,
	mailer *mail.Dialer,
	captchaManager *captchas.Manager,
	accessManager *access.Manager,
) *api.Application {
	app := newApp(logger, db, view, sessionManager, userManager, mailer, captchaManager, accessManager)
	return &api.Application{Application: app}
}

type apiRoutes routeutil.Routes

func provideAPIRoutes(
	accessManager *access.Manager,
	post *controllers.Post,
	user *controllers.User,
) apiRoutes {
	return apiRoutes{
		routeutil.NewRoute(http.MethodGet, "/login", user.Login).Name("login"),
		routeutil.NewRoute(http.MethodGet, "/logout", user.Logout).Name("logout"),
		routeutil.NewRoute(http.MethodGet, "/users", user.Index).Name("users"),
		routeutil.NewRoute(http.MethodGet, "/users/:id", user.Index).Name("user"),

		routeutil.NewRoute(http.MethodGet, "/posts", post.Index).Name("posts"),
		routeutil.NewRoute(http.MethodGet, "/posts/:id", post.View),
		routeutil.NewRoute(http.MethodPost, "/posts/:id", post.Create),
		routeutil.NewRoute(http.MethodPut, "/posts/:id", post.Update),
		routeutil.NewRoute(http.MethodDelete, "/posts/:id", post.Delete),
	}
}
