package cmd

import (
	"net/http"

	"github.com/clevergo/demo/internal/backend/controllers"
	"github.com/clevergo/demo/pkg/access"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/google/wire"
)

var backendAppSet = wire.NewSet(
	provideBackendApp, provideBackendView,
	provideBackendRoutes,
	controllers.NewSite, controllers.NewPost,
)

type backendRoutes routeutil.Routes

func provideBackendRoutes(
	accessManager *access.Manager,
	site *controllers.Site,
	post *controllers.Post,
) backendRoutes {
	return backendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),

		routeutil.NewRoute(http.MethodGet, "/post", post.Index).Name("post").Middlewares(
			accessManager.Middleware("post", "read"),
		),
	}
}
