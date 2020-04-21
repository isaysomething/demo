package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/authz"
)

type AuthzMiddleware clevergo.MiddlewareFunc

func NewAuthzMiddleware(enforcer *casbin.Enforcer, userManager *UserManager) AuthzMiddleware {
	authorization := authz.New(enforcer, userManager.Manager)
	authorization.Skipper = clevergo.PathSkipper(
		"/v1/captcha",
		"/v1/check-captcha",
		"/v1/user/login",
	)

	return AuthzMiddleware(authorization.Middleware())
}
