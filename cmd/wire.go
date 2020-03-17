//+build wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	provideRouter, provideMiddlewares, provideI18N,
	provideLogger, provideDB, provideSessionManager, provideSessionStore, provideUserManager,
	provideIdentityStore, provideMailer, provideCaptchaManager,
	provideEnforcer, access.New,
	provideTencentClient, provideTencentCaptcha,
)

func initializeServer() (*web.Server, func(), error) {
	wire.Build(superSet, frontendSet, backendAppSet, provideServer)
	return &web.Server{}, nil, nil
}

func initializeAPIServer() (*web.Server, func(), error) {
	wire.Build(superSet, apiSet, provideAPIServer)
	return &web.Server{}, nil, nil
}
