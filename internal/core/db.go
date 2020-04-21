package core

import (
	"github.com/clevergo/demo/pkg/db"
)

// DBConfig is a database config.
type DBConfig struct {
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
}

// NewDB returns a sqlx.DB with the given config.
func NewDB(cfg DBConfig) (*db.DB, func(), error) {
	conn, err := db.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, nil, err
	}
	conn.SetLogger(&db.StdLogger{})

	return conn, func() {
		if err := conn.Close(); err != nil {
		}
	}, nil
}
