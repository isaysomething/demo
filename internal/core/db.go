package core

import "github.com/jmoiron/sqlx"

// DBConfig is a database config.
type DBConfig struct {
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
}

// NewDB returns a sqlx.DB with the given config.
func NewDB(cfg DBConfig) (*sqlx.DB, error) {
	return sqlx.Open(cfg.Driver, cfg.DSN)
}
