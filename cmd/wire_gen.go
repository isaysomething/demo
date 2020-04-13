// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Injectors from wire.go:

func initializeServer() (*core.Server, func(), error) {
	logConfig := provideLogConfig()
	logger, cleanup, err := core.NewLogger(logConfig)
	if err != nil {
		return nil, nil, err
	}
	params := provideParams()
	dbConfig := provideDBConfig()
	db, cleanup2, err := core.NewDB(dbConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	viewConfig := provideViewConfig()
	viewManager := core.NewViewManager(viewConfig)
	sessionConfig := provideSessionConfig()
	redisConfig := provideRedisConfig()
	store := core.NewSessionStore(redisConfig)
	sessionManager := core.NewSessionManager(sessionConfig, store)
	jwtConfig := provideJWTConfig()
	jwtManager := core.NewJWTManager(jwtConfig)
	identityStore := core.NewIdentityStore(db, jwtManager)
	manager := core.NewUserManager(identityStore, sessionManager)
	mailerConfig := provideMailerConfig()
	dialer := core.NewMailer(mailerConfig)
	enforcer, err := core.NewEnforcer(dbConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	accessManager := access.New(enforcer, manager)
	application := frontend.New(logger, params, db, viewManager, sessionManager, manager, dialer, accessManager)
	captchaConfig := provideCaptchaConfig()
	captchasStore := core.NewCaptchaStore(redisConfig)
	captchasManager := core.NewCaptchaManager(captchaConfig, captchasStore)
	csrfConfig := provideCSRFConfig()
	csrfMiddleware := core.NewCSRFMiddleware(csrfConfig)
	i18NConfig := provideI18NConfig()
	i18nStore := core.NewFileStore(i18NConfig)
	translators, err := core.NewI18N(i18NConfig, i18nStore)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	v := core.NewI18NLanguageParsers(i18NConfig)
	i18NMiddleware := core.NewI18NMiddleware(translators, v)
	gzipMiddleware := core.NewGzipMiddleware()
	sessionMiddleware := core.NewSessionMiddleware(sessionManager)
	m := core.NewMinify()
	minifyMiddleware := core.NewMinifyMiddleware(m)
	loggingMiddleware := core.NewLoggingMiddleware()
	router := provideRouter(application, captchasManager, csrfMiddleware, i18NMiddleware, gzipMiddleware, sessionMiddleware, minifyMiddleware, loggingMiddleware)
	server := provideServer(router, logger)
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

func initializeAPIServer() (*core.Server, func(), error) {
	logConfig := provideLogConfig()
	logger, cleanup, err := core.NewLogger(logConfig)
	if err != nil {
		return nil, nil, err
	}
	params := provideParams()
	dbConfig := provideDBConfig()
	db, cleanup2, err := core.NewDB(dbConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sessionConfig := provideSessionConfig()
	redisConfig := provideRedisConfig()
	store := core.NewSessionStore(redisConfig)
	sessionManager := core.NewSessionManager(sessionConfig, store)
	jwtConfig := provideJWTConfig()
	jwtManager := core.NewJWTManager(jwtConfig)
	identityStore := core.NewIdentityStore(db, jwtManager)
	userManager := api.NewUserManager(identityStore)
	mailerConfig := provideMailerConfig()
	dialer := core.NewMailer(mailerConfig)
	enforcer, err := core.NewEnforcer(dbConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	manager := core.NewUserManager(identityStore, sessionManager)
	accessManager := access.New(enforcer, manager)
	application := api.New(logger, params, db, sessionManager, userManager, dialer, accessManager)
	captchaConfig := provideCaptchaConfig()
	captchasStore := core.NewCaptchaStore(redisConfig)
	captchasManager := core.NewCaptchaManager(captchaConfig, captchasStore)
	authenticator := core.NewAuthenticator(identityStore)
	corsConfig := provideCORSConfig()
	corsMiddleware := core.NewCORSMiddleware(corsConfig)
	server := provideAPIServer(logger, application, captchasManager, jwtManager, userManager, authenticator, corsMiddleware)
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var superSet = wire.NewSet(
	configSet, core.NewDB, core.NewSessionStore, core.NewSessionManager, core.NewMailer, core.NewLogger, core.NewAuthenticator, core.NewIdentityStore, core.NewUserManager, core.NewCaptchaStore, core.NewCaptchaManager, core.NewI18N, core.NewFileStore, core.NewI18NLanguageParsers, provideRouter, core.NewEnforcer, access.New, core.NewJWTManager, core.MiddlewareSet,
)
