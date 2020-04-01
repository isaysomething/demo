package core

import "github.com/google/wire"

var MiddlewareSet = wire.NewSet(
	NewCSRFMiddleware,
	NewI18NMiddleware,
	NewCORSMiddleware,
	NewGzipMiddleware,
	NewSessionMiddleware,
	NewMinify, NewMinifyMiddleware,
	NewLoggingMiddleware,
)
