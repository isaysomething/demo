package middlewares

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/gorilla/csrf"
)

func CSRF() clevergo.MiddlewareFunc {
	m := csrf.Protect([]byte("123456"), csrf.Secure(false))
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) (err error) {
			handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx.Response = w
				ctx.Request = r
				err = next(ctx)
			}))
			handler.ServeHTTP(ctx.Response, ctx.Request)
			return
		}
	}
}
