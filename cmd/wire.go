//+build wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/web"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	provideServer, provideRouter, provideMiddlewares, provideI18N,
	provideLogger, provideDB, provideSessionManager, provideSessionStore, provideUserStore,
	provideIdentityStore, provideMailer, provideCaptchaManager,
)

func initializeServer() (*web.Server, func(), error) {
	wire.Build(superSet, frontendSet, backendAppSet)
	return &web.Server{}, nil, nil
}
