package minify

import (
	"net/http"

	"github.com/tdewolff/minify/v2"
)

func Minify(next http.Handler) http.Handler {
	m := minify.New()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(m.ResponseWriter(w, r), r)
	})
}
