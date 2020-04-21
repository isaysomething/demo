package db

import "log"

type Logger interface {
	Log(query string, args ...interface{})
}

type StdLogger struct {
}

func (l *StdLogger) Log(query string, args ...interface{}) {
	log.Printf("SQL: %s, args: %+v\n", query, args)
}
