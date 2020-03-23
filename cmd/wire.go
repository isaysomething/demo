//+build wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	provideRouter, provideMiddlewares, provideI18N,
	provideLogger, provideDB, provideSessionManager, provideSessionStore, provideUserManager,
	provideIdentityStore, provideMailer, provideCaptchaManager,
	provideEnforcer, access.New, provideAuthenticator,
)

func initializeServer() (*core.Server, func(), error) {
	wire.Build(superSet, frontendSet, provideServer)
	return &core.Server{}, nil, nil
}

func initializeAPIServer() (*core.Server, func(), error) {
	wire.Build(superSet, apiSet, provideAPIServer)
	return &core.Server{}, nil, nil
}
