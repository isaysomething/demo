package middlewares

import (
	"net/http"
	"strings"
)

type Skipper func(r *http.Request) bool

func NewPathSkipper(paths ...string) Skipper {
	return func(r *http.Request) bool {
		path := r.URL.Path
		for _, p := range paths {
			if p == path {
				return true
			}

			if p[len(p)-1] == '*' && strings.HasPrefix(path, p[:len(p)-1]) {
				return true
			}
		}

		return false
	}
}
