package middlewares

import "github.com/clevergo/clevergo"

func ServerHeader(v string) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) error {
			ctx.Response.Header().Set("Server", v)
			return next(ctx)
		}
	}
}
