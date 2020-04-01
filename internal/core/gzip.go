package core

import (
	"github.com/NYTimes/gziphandler"
	"github.com/clevergo/clevergo"
)

type GzipMiddleware clevergo.MiddlewareFunc

func NewGzipMiddleware() GzipMiddleware {
	return GzipMiddleware(clevergo.WrapHH(gziphandler.GzipHandler))
}
