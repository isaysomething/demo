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
	user *controllers.User,
) frontendRoutes {
	return frontendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),
		routeutil.NewRoute(http.MethodGet, "/about", site.About).Name("about"),
		routeutil.NewRoute(http.MethodGet, "/contact", site.Contact).Name("contact"),
		routeutil.NewRoute(http.MethodPost, "/contact", site.Contact),

		routeutil.NewRoute(http.MethodGet, "/login", user.Login).Name("login"),
		routeutil.NewRoute(http.MethodPost, "/login", user.Login),
		routeutil.NewRoute(http.MethodPost, "/logout", user.Logout).Name("logout"),
		routeutil.NewRoute(http.MethodGet, "/sign-up", user.SignUp).Name("sign-up"),
		routeutil.NewRoute(http.MethodPost, "/sign-up", user.SignUp),
		routeutil.NewRoute(http.MethodGet, "/verify-email", user.VerifyEmail).Name("verify-email"),
		routeutil.NewRoute(http.MethodGet, "/resend-verification-email", user.ResendVerificationEmail).Name("resend-verification-email"),
		routeutil.NewRoute(http.MethodPost, "/resend-verification-email", user.ResendVerificationEmail),
		routeutil.NewRoute(http.MethodGet, "/reset-password", user.ResetPassword).Name("reset-password"),
		routeutil.NewRoute(http.MethodPost, "/reset-password", user.ResetPassword),
		routeutil.NewRoute(http.MethodGet, "/change-password", user.ChangePassword).Name("change-password"),
		routeutil.NewRoute(http.MethodPost, "/change-password", user.ChangePassword),
		routeutil.NewRoute(http.MethodGet, "/setting", user.Setting).Name("setting"),
	}
}
