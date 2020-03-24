package authz

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/users"
)

func New(enforcer *casbin.Enforcer, userManager *users.Manager) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) error {
			user, _ := userManager.Get(ctx.Request, ctx.Response)
			if user.IsGuest() {
				ctx.Error(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return nil
			}
			ok, err := enforcer.Enforce(user.GetIdentity().GetID, ctx.Request.URL.Path, ctx.Request.Method)
			if err != nil {
				log.Println(err)
			}
			if !ok {
				ctx.Error(http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return nil
			}
			return next(ctx)
		}
	}
}
