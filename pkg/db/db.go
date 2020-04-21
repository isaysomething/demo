package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	logger Logger
}

func Open(driver, dsn string) (*DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB: db,
	}, nil
}

func MustOpen(driver, dsn string) *DB {
	return &DB{
		DB: sqlx.MustOpen(driver, dsn),
	}
}

func (db *DB) SetLogger(logger Logger) {
	db.logger = logger
}

func (db *DB) log(query string, args ...interface{}) {
	if db.logger != nil {
		db.logger.Log(query, args...)
	}
}

func (db *DB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	db.log(query, arg)
	return db.DB.NamedQuery(query, arg)
}

func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	db.log(query, arg)
	return db.DB.NamedExec(query, arg)
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	db.log(query, args...)
	return db.DB.Get(dest, query, args...)
}

func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	db.log(query, args...)
	return db.DB.Select(dest, query, args...)
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db.log(query, args...)
	return db.DB.Query(query, args...)
}

func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	db.log(query, args...)
	return db.DB.Queryx(query, args...)
}

func (db *DB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	db.log(query, args...)
	return db.DB.QueryRowx(query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.log(query, args...)
	return db.DB.Exec(query, args...)
}
