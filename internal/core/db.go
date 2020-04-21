package core

import (
	"github.com/clevergo/demo/pkg/sqlex"
)

// DBConfig is a database config.
type DBConfig struct {
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
}

// NewDB returns a sqlx.DB with the given config.
func NewDB(cfg DBConfig) (*sqlex.DB, func(), error) {
	db, err := sqlex.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, nil, err
	}
	db.SetLogger(&sqlex.StdLogger{})

	return db, func() {
		if err := db.Close(); err != nil {
		}
	}, nil
}
