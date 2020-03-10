package cmd

import (
	"net/http"

	"github.com/clevergo/demo/internal/backend/controllers"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/google/wire"
)

var backendAppSet = wire.NewSet(
	provideBackendApp, provideAdminView,
	provideBackendRoutes,
	controllers.NewSite, controllers.NewPost,
)

type backendRoutes routeutil.Routes

func provideBackendRoutes(
	site *controllers.Site,
	post *controllers.Post,
) backendRoutes {
	return backendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),
	}
}
