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
	controllers.NewSite, controllers.NewUser, controllers.NewPost,
)

type backendRoutes routeutil.Routes

func provideBackendRoutes(
	accessManager *access.Manager,
	site *controllers.Site,
	post *controllers.Post,
	user *controllers.User,
) backendRoutes {
	return backendRoutes{
		routeutil.NewRoute(http.MethodGet, "/", site.Index).Name("home"),
		routeutil.NewRoute(http.MethodPost, "/captcha", site.Captcha).Name("captcha"),
		routeutil.NewRoute(http.MethodPost, "/check-captcha", site.CheckCaptcha),
		routeutil.NewRoute(http.MethodPost, "/user/check-username", user.CheckUsername),
		routeutil.NewRoute(http.MethodPost, "/user/check-email", user.CheckEmail),

		routeutil.NewRoute(http.MethodGet, "/post", post.Index).Name("post").Middlewares(
			accessManager.Middleware("post", "read"),
		),
		routeutil.NewRoute(http.MethodGet, "/post/create", post.Create).Name("post-create").Middlewares(
			accessManager.Middleware("post", "read"),
		),
		routeutil.NewRoute(http.MethodPost, "/post/create", post.Create).Middlewares(
			accessManager.Middleware("post", "read"),
		),

		routeutil.NewRoute(http.MethodGet, "/users/:user/:role", user.Index).Name("user"),
		routeutil.NewRoute(http.MethodGet, "/login", user.Login).Name("login"),
		routeutil.NewRoute(http.MethodPost, "/login", user.Login),
		routeutil.NewRoute(http.MethodPost, "/logout", user.Logout).Name("logout"),
		routeutil.NewRoute(http.MethodGet, "/signup", user.Signup).Name("signup"),
		routeutil.NewRoute(http.MethodPost, "/signup", user.Signup),
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
