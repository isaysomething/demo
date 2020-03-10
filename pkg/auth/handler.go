package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/clevergo/auth"
)

type contextKey int

const (
	identityKey contextKey = iota
)

func GetIdentity(r *http.Request) auth.Identity {
	v, _ := r.Context().Value(identityKey).(auth.Identity)
	return v
}

type Option func(*Handler)

func IsOptionalFunc(f func(*http.Request) bool) Option {
	return func(h *Handler) {
		h.isOptional = f
	}
}

func OptionalPaths(paths ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		path := r.URL.Path
		for _, v := range paths {
			if path == v {
				return true
			}
			if v[len(v)-1] == '*' && strings.HasPrefix(path, v[:]) {
				return true
			}
		}

		return false
	}
}

type Handler struct {
	authenticator auth.Authenticator
	next          http.Handler
	isOptional    func(*http.Request) bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	identity, err := h.authenticator.Authenticate(r)
	if err != nil && (h.isOptional != nil && !h.isOptional(r)) {
		h.authenticator.Challenge(w)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), identityKey, identity)
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func NewHandler(authenticator auth.Authenticator, next http.Handler, opts ...Option) *Handler {
	h := &Handler{authenticator: authenticator, next: next}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func Middleware(authenticator auth.Authenticator, opts ...Option) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewHandler(authenticator, next, opts...)
	}
}
