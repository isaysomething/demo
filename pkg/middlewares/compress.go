package middlewares

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/gorilla/handlers"
)

func Wrap(h http.Handler) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) (err error) {
			h.ServeHTTP(ctx.Response, ctx.Request)
			return next(ctx)
		}
	}
}

func Compress(level int) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) (err error) {
			handler := handlers.CompressHandlerLevel(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx.Response = w
				ctx.Request = r
				err = next(ctx)
			}), level)
			handler.ServeHTTP(ctx.Response, ctx.Request)
			return
		}
	}
}
