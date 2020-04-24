package sqlex

import (
	"log"
	"time"
)

type Logger interface {
	Log(duration time.Duration, query string, args ...interface{})
}

var StdLogger Logger = &logger{}

type logger struct {
}

func (l *logger) Log(duration time.Duration, query string, args ...interface{}) {
	log.Printf("[%s] [%s] %v\n", duration.String(), query, args)
}
