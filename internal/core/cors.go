package core

import (
	"github.com/clevergo/clevergo"
	"github.com/rs/cors"
)

type CORSConfig struct {
	AllowedOrigins     []string `koanf:"allowed_origins"`
	AllowedMethods     []string `koanf:"allowed_methods"`
	AllowedHeaders     []string `koanf:"allowed_headers"`
	MaxAge             int      `koanf:"max_age"`
	AllowedCredentials bool     `koanf:"allow_credentials"`
	Debug              bool     `koanf:"debug"`
}

func NewCORS(cfg CORSConfig) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedHeaders:   cfg.AllowedHeaders,
		MaxAge:           cfg.MaxAge,
		AllowCredentials: cfg.AllowedCredentials,
		AllowedMethods:   cfg.AllowedMethods,
		Debug:            cfg.Debug,
	})
}

type CORSMiddleware clevergo.MiddlewareFunc

func NewCORSMiddleware(cfg CORSConfig) CORSMiddleware {
	fn := NewCORS(cfg).Handler
	return CORSMiddleware(clevergo.WrapHH(fn))
}
