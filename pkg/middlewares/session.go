package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/clevergo"
)

func Session(manager *scs.SessionManager) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) (err error) {
			handler := manager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx.Response = w
				ctx.Request = r
				err = next(ctx)
			}))
			handler.ServeHTTP(ctx.Response, ctx.Request)
			return
		}
	}
}
