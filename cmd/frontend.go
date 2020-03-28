package cmd

import (
	"net/http"

	commonctl "github.com/clevergo/demo/internal/controllers"
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
	captcha *commonctl.Captcha,
	user *controllers.User,
) frontendRoutes {
	return frontendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),
		routeutil.NewRoute(http.MethodGet, "/robots.txt", site.Robots),
		routeutil.NewRoute(http.MethodGet, "/about", site.About).Name("about"),
		routeutil.NewRoute(http.MethodGet, "/contact", site.Contact).Name("contact"),
		routeutil.NewRoute(http.MethodPost, "/contact", site.Contact),
		routeutil.NewRoute(http.MethodPost, "/captcha", captcha.Generate).Name("captcha"),
		routeutil.NewRoute(http.MethodPost, "/check-captcha", captcha.Verify),

		routeutil.NewRoute(http.MethodGet, "/user", user.Index).Name("user"),
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
		routeutil.NewRoute(http.MethodGet, "/user/verify-email", user.VerifyEmail).Name("verify-email"),
		routeutil.NewRoute(http.MethodGet, "/user/resend-verification-email", user.ResendVerificationEmail).Name("resend-verification-email"),
		routeutil.NewRoute(http.MethodPost, "/user/resend-verification-email", user.ResendVerificationEmail),
		routeutil.NewRoute(http.MethodPost, "/user/change-password", user.ChangePassword).Name("change-password"),
	}
}
