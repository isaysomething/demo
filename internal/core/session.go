package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/sessionmidware"
	"github.com/gomodule/redigo/redis"
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

func NewSessionManager(cfg SessionConfig, store scs.Store) *scs.SessionManager {
	m := scs.New()
	m.Store = store
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

func NewSessionStore(cfg RedisConfig) scs.Store {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	opts := []redis.DialOption{
		redis.DialDatabase(cfg.Database),
	}
	if cfg.Password != "" {
		opts = append(opts, redis.DialPassword(cfg.Password))
	}
	return redisstore.New(redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", address, opts...)
	}, 1000))
}

type SessionMiddleware clevergo.MiddlewareFunc

func NewSessionMiddleware(manager *scs.SessionManager) SessionMiddleware {
	return SessionMiddleware(sessionmidware.New(manager))
}
