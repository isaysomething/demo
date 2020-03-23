package middlewares

import (
	"regexp"

	"github.com/clevergo/clevergo"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

func Minify() clevergo.MiddlewareFunc {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) error {
			w := m.ResponseWriter(ctx.Response, ctx.Request)
			defer w.Close()
			ctx.Response = w
			return next(ctx)
		}
	}
}
