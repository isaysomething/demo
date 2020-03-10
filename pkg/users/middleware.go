package users

import (
	"net/http"

	"github.com/clevergo/auth"
)

func Handler(next http.Handler, store *Store, authenticator auth.Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := store.Get(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user.IsGuest() {
			identity, err := authenticator.Authenticate(r)
			if err == nil {
				user.Login(r, w, identity, 0)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func Middleware(store *Store, authenticator auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Handler(next, store, authenticator)
	}
}
