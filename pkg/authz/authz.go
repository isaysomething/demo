package authz

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/users"
)

var (
	ErrUnauthorized = clevergo.NewError(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	ErrForbidden    = clevergo.NewError(http.StatusForbidden, errors.New("you are not allowed to access this page"))
)

type Authorization struct {
	enforcer    *casbin.Enforcer
	userManager *users.Manager
	Skipper     clevergo.Skipper
}

func New(enforcer *casbin.Enforcer, userManager *users.Manager) *Authorization {
	return &Authorization{enforcer: enforcer, userManager: userManager}
}

func (a *Authorization) Middleware() clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) error {
			if a.Skipper == nil || !a.Skipper(ctx) {
				user, _ := a.userManager.Get(ctx.Request, ctx.Response)
				if user.IsGuest() {
					return ErrUnauthorized
				}
				fmt.Println(a.getSubject(user), a.getObject(ctx), a.getAction(ctx))
				ok, err := a.enforcer.Enforce(a.getSubject(user), a.getObject(ctx), a.getAction(ctx))
				if err != nil {
					return clevergo.NewError(http.StatusInternalServerError, err)
				}
				if !ok {
					return ErrForbidden
				}
			}
			return next(ctx)
		}
	}
}

func (a *Authorization) getSubject(user *users.User) interface{} {
	return "user_" + user.GetIdentity().GetID()
}

func (a *Authorization) getObject(ctx *clevergo.Context) interface{} {
	return ctx.Request.URL.Path
}

func (a *Authorization) getAction(ctx *clevergo.Context) interface{} {
	return ctx.Request.Method
}
