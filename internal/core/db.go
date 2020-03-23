package core

import "github.com/jmoiron/sqlx"

// DBConfig is a database config.
type DBConfig struct {
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
}

// NewDB returns a sqlx.DB with the given config.
func NewDB(cfg DBConfig) (*sqlx.DB, func(), error) {
	db, err := sqlx.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
		}
	}, nil
}
