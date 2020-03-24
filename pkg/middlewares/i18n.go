package middlewares

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/i18n"
)

func I18N(translators *i18n.Translators, parsers ...i18n.LanguageParser) clevergo.MiddlewareFunc {
	m := i18n.Middleware(translators, parsers...)
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
