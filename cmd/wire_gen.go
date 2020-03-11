// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	controllers2 "github.com/clevergo/demo/internal/backend/controllers"
	"github.com/clevergo/demo/internal/frontend/controllers"
	"github.com/clevergo/demo/internal/web"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Injectors from wire.go:

func initializeServer() (*web.Server, func(), error) {
	logger, cleanup, err := provideLogger()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := provideDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	manager := provideView()
	store := provideSessionStore()
	sessionManager := provideSessionManager(store)
	identityStore := provideIdentityStore(db)
	usersManager := provideUserManager(identityStore, sessionManager)
	dialer := provideMailer()
	captchasManager := provideCaptchaManager()
	application := provideApp(logger, db, manager, sessionManager, usersManager, dialer, captchasManager)
	site := controllers.NewSite(application)
	user := controllers.NewUser(application)
	cmdFrontendRoutes := provideFrontendRoutes(site, user)
	adminView := provideAdminView(logger)
	backendApplication := provideBackendApp(logger, db, adminView, sessionManager, usersManager, dialer, captchasManager)
	controllersSite := controllers2.NewSite(backendApplication)
	post := controllers2.NewPost(backendApplication)
	cmdBackendRoutes := provideBackendRoutes(controllersSite, post)
	router := provideRouter(application, cmdFrontendRoutes, backendApplication, cmdBackendRoutes)
	translators, err := provideI18N()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	v, err := provideMiddlewares(sessionManager, translators, usersManager)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	server := provideServer(router, logger, v)
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var superSet = wire.NewSet(
	provideServer, provideRouter, provideMiddlewares, provideI18N,
	provideLogger, provideDB, provideSessionManager, provideSessionStore, provideUserManager,
	provideIdentityStore, provideMailer, provideCaptchaManager,
)
