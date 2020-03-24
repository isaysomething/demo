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
		routeutil.NewRoute(http.MethodGet, "/robots.txt", site.Robots),
		routeutil.NewRoute(http.MethodGet, "/about", site.About).Name("about"),
		routeutil.NewRoute(http.MethodGet, "/contact", site.Contact).Name("contact"),
		routeutil.NewRoute(http.MethodPost, "/contact", site.Contact),
		routeutil.NewRoute(http.MethodPost, "/captcha", site.Captcha).Name("captcha"),
		routeutil.NewRoute(http.MethodPost, "/check-captcha", site.CheckCaptcha),

		routeutil.NewRoute(http.MethodGet, "/login", user.Login).Name("login"),
		routeutil.NewRoute(http.MethodPost, "/login", user.Login),
		routeutil.NewRoute(http.MethodPost, "/logout", user.Logout).Name("logout"),
		routeutil.NewRoute(http.MethodGet, "/signup", user.Signup).Name("signup"),
		routeutil.NewRoute(http.MethodPost, "/signup", user.Signup),
		routeutil.NewRoute(http.MethodPost, "/user/check-email", user.CheckEmail),
		routeutil.NewRoute(http.MethodPost, "/user/check-username", user.CheckUsername),
		routeutil.NewRoute(http.MethodGet, "/user/request-reset-password", user.RequestResetPassword).Name("request-reset-password"),
		routeutil.NewRoute(http.MethodPost, "/user/request-reset-password", user.RequestResetPassword),
		routeutil.NewRoute(http.MethodGet, "/user/reset-password", user.ResetPassword).Name("reset-password"),
		routeutil.NewRoute(http.MethodPost, "/user/reset-password", user.ResetPassword),

		routeutil.NewRoute(http.MethodGet, "/users/:user/:role", user.Index).Name("user"),
		routeutil.NewRoute(http.MethodGet, "/verify-email", user.VerifyEmail).Name("verify-email"),
		routeutil.NewRoute(http.MethodGet, "/resend-verification-email", user.ResendVerificationEmail).Name("resend-verification-email"),
		routeutil.NewRoute(http.MethodPost, "/resend-verification-email", user.ResendVerificationEmail),
		routeutil.NewRoute(http.MethodGet, "/change-password", user.ChangePassword).Name("change-password"),
		routeutil.NewRoute(http.MethodPost, "/change-password", user.ChangePassword),
		routeutil.NewRoute(http.MethodGet, "/setting", user.Setting).Name("setting"),
	}
}
