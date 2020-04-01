package core

import (
	"github.com/clevergo/clevergo"
	"github.com/gorilla/csrf"
)

type CSRFConfig struct {
	AuthKey        string   `koanf:"auth_key"`
	MaxAge         int      `koanf:"max_age"`
	Domain         string   `koanf:"domain"`
	Path           string   `koanf:"path"`
	Secure         bool     `koanf:"secure"`
	HTTPOnly       bool     `koanf:"http_only"`
	SameSite       int      `koanf:"same_site"`
	RequestHeader  string   `koanf:"request_header"`
	FieldName      string   `koanf:"field_name"`
	CookieName     string   `koanf:"cookie_name"`
	TrustedOrigins []string `koanf:"trusted_origins"`
}

type CSRFMiddleware clevergo.MiddlewareFunc

func NewCSRFMiddleware(cfg CSRFConfig) CSRFMiddleware {
	opts := []csrf.Option{
		csrf.MaxAge(cfg.MaxAge),
		csrf.Path(cfg.Path),
		csrf.Secure(cfg.Secure),
		csrf.HttpOnly(cfg.HTTPOnly),
		csrf.SameSite(csrf.SameSiteMode(cfg.SameSite)),
	}
	if cfg.RequestHeader != "" {
		opts = append(opts, csrf.RequestHeader(cfg.RequestHeader))
	}
	if cfg.FieldName != "" {
		opts = append(opts, csrf.FieldName(cfg.FieldName))
	}
	if cfg.CookieName != "" {
		opts = append(opts, csrf.CookieName(cfg.CookieName))
	}
	if len(cfg.TrustedOrigins) > 0 {
		opts = append(opts, csrf.TrustedOrigins(cfg.TrustedOrigins))
	}

	fn := csrf.Protect([]byte(cfg.AuthKey), opts...)
	return CSRFMiddleware(clevergo.WrapHH(fn))
}
