//+build wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	provideServer, provideRouter, provideMiddlewares, provideI18N,
	provideLogger, provideDB, provideSessionManager, provideSessionStore, provideUserManager,
	provideIdentityStore, provideMailer, provideCaptchaManager,
	provideEnforcer, access.New,
)

func initializeServer() (*web.Server, func(), error) {
	wire.Build(superSet, frontendSet, backendAppSet, apiSet)
	return &web.Server{}, nil, nil
}
