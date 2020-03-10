package middlewares

import "net/http"

type LoginChecker struct {
	loginUrl string
	skipper  Skipper
	isGuest  func(r *http.Request, w http.ResponseWriter) bool
	handler  http.Handler
}

func NewLoginChecker(handler http.Handler, isGuest func(r *http.Request, w http.ResponseWriter) bool, skipper Skipper) *LoginChecker {
	return &LoginChecker{
		loginUrl: "/admin/login",
		isGuest:  isGuest,
		skipper:  skipper,
		handler:  handler,
	}
}

func LoginCheckerMiddleware(isGuest func(r *http.Request, w http.ResponseWriter) bool, skipper Skipper) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewLoginChecker(next, isGuest, skipper)
	}
}

func (lc *LoginChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if lc.isGuest(r, w) && !lc.skipper(r) {
		http.Redirect(w, r, lc.loginUrl, http.StatusFound)
		return
	}

	lc.handler.ServeHTTP(w, r)
}
