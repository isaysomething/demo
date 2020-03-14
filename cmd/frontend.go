package cmd

import (
	"net/http"

	"github.com/clevergo/demo/internal/frontend/controllers"
	"github.com/clevergo/demo/pkg/routeutil"
	"github.com/google/wire"
)

var frontendSet = wire.NewSet(
	provideApp, provideView,
	provideFrontendRoutes,
	controllers.NewSite, controllers.NewUser,
)

type frontendRoutes routeutil.Routes

func provideFrontendRoutes(
	site *controllers.Site,
) frontendRoutes {
	return frontendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),
		routeutil.NewRoute(http.MethodGet, "/about", site.About).Name("about"),
		routeutil.NewRoute(http.MethodGet, "/contact", site.Contact).Name("contact"),
		routeutil.NewRoute(http.MethodPost, "/contact", site.Contact),
	}
}
