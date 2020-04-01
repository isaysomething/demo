package core

import (
	"os"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/middleware"
)

type LoggingMiddleware clevergo.MiddlewareFunc

func NewLoggingMiddleware() LoggingMiddleware {
	return LoggingMiddleware(clevergo.WrapHH(middleware.Logging(os.Stdout)))
}
