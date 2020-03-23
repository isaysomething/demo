package middlewares

import (
	"io"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/gorilla/handlers"
)

func Logging(wr io.Writer) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) (err error) {
			handler := handlers.LoggingHandler(wr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx.Response = w
				ctx.Request = r
				err = next(ctx)
			}))
			handler.ServeHTTP(ctx.Response, ctx.Request)
			return
		}
	}
}
