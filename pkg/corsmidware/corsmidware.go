package corsmidware

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/rs/cors"
)

type corsWrapper struct {
	*cors.Cors
	optionPassthrough bool
}

func (c *corsWrapper) middleware(next clevergo.Handle) clevergo.Handle {
	return func(ctx *clevergo.Context) error {
		c.HandlerFunc(ctx.Response, ctx.Request)
		if !c.optionPassthrough &&
			ctx.Request.Method == http.MethodOptions &&
			ctx.Request.Header.Get("Access-Control-Request-Method") != "" {
			return nil
		}
		return next(ctx)
	}
}

func Default() clevergo.MiddlewareFunc {
	c := &corsWrapper{Cors: cors.Default()}
	return c.middleware
}

func New(options cors.Options) clevergo.MiddlewareFunc {
	c := &corsWrapper{Cors: cors.New(options), optionPassthrough: options.OptionsPassthrough}
	return c.middleware
}
