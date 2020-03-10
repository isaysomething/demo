package web

import "github.com/jmoiron/sqlx"

func NewDB(cfg DBConfig) (*sqlx.DB, error) {
	//return sql.Open(cfg.Driver, cfg.DSN)
	return sqlx.Open(cfg.Driver, cfg.DSN)
}
