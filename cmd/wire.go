//+build wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/controllers"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	configSet,
	core.NewDB,
	core.NewSessionStore, core.NewSessionManager,
	core.NewMailer,
	core.NewLogger,
	core.NewAuthenticator, core.NewIdentityStore, core.NewUserManager,
	core.NewCaptchaStore, core.NewCaptchaManager,
	core.NewI18N, core.NewFileStore, core.NewI18NLanguageParsers,
	provideRouter,
	provideEnforcer, access.New,

	// controllers
	controllers.NewCaptcha,
)

func initializeServer() (*core.Server, func(), error) {
	wire.Build(superSet, frontendSet, provideServer)
	return &core.Server{}, nil, nil
}

func initializeAPIServer() (*core.Server, func(), error) {
	wire.Build(superSet, apiSet, provideAPIServer)
	return &core.Server{}, nil, nil
}
