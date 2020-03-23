package core

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

type SessionConfig struct {
	Lifetime       int    `koanf:"lifetime"`
	IdleTimeout    int    `koanf:"idle_timeout"`
	CookieName     string `koanf:"cookie_name"`
	CookieDomain   string `koanf:"cookie_domain"`
	CookiePath     string `koanf:"cookie_path"`
	CookieHTTPOnly bool   `koanf:"cookie_http_only"`
	CookiePersist  bool   `koanf:"cookie_persist"`
	CookieSecure   bool   `koanf:"cookie_secure"`
	CookieSameSite int    `koanf:"cookie_same_site"`
}

func NewSessionManager(cfg SessionConfig) *scs.SessionManager {
	m := scs.New()
	m.Lifetime = time.Duration(cfg.Lifetime) * time.Second
	m.IdleTimeout = time.Duration(cfg.IdleTimeout) * time.Second
	m.Cookie.Name = cfg.CookieName
	m.Cookie.Domain = cfg.CookieDomain
	m.Cookie.HttpOnly = cfg.CookieHTTPOnly
	m.Cookie.Path = cfg.CookiePath
	m.Cookie.Persist = cfg.CookiePersist
	m.Cookie.SameSite = http.SameSite(cfg.CookieSameSite)
	m.Cookie.Secure = cfg.CookieSecure
	return m
}
