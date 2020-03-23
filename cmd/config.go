package cmd

import (
	"io/ioutil"

	"github.com/clevergo/demo/internal/core"
	"github.com/gobuffalo/packr/v2"
	"github.com/google/wire"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
)

var configSet = wire.NewSet(
	provideDBConfig,
	provideRedisConfig,
	provideSessionConfig,
	provideMailerConfig,
	provideLogConfig,
	provideCaptchaConfig,
)

func provideDBConfig() core.DBConfig {
	return cfg.DB
}

func provideMailerConfig() core.MailerConfig {
	return cfg.Mailer
}

func provideRedisConfig() core.RedisConfig {
	return cfg.Redis
}

func provideSessionConfig() core.SessionConfig {
	return cfg.Session
}

func provideLogConfig() core.LogConfig {
	return cfg.Log
}

func provideCaptchaConfig() core.CaptchaConfig {
	return cfg.Captcha
}

func parseConfig() error {
	parser := toml.Parser()
	configFS := packr.New("configs", "./../configs")
	// load default configurations.
	configs := configFS.List()
	for _, name := range configs {
		f, err := configFS.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if err := k.Load(rawbytes.Provider(content), parser); err != nil {
			return err
		}
	}

	if err := k.Load(file.Provider(*cfgFile), parser); err != nil {
		return err
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		return err
	}

	return nil
}
