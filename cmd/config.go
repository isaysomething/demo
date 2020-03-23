package cmd

import (
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/pkg/params"
)

type Config struct {
	Server struct {
		Addr string `koanf:"addr"`
		Root string `koanf:"root"`

		SSL         bool   `koanf:"ssl"`
		SSLCertFile string `koanf:"ssl_cert_file"`
		SSLKeyFile  string `koanf:"ssl_key_file"`

		Log core.LogConfig `koanf:"log"`

		AccessLog         bool   `koanf:"access_log"`
		AccessLogFile     string `koanf:"access_log_file"`
		AccessLogFileMode uint32 `koanf:"access_log_file_mode"`

		Gzip      bool `koanf:"gzip"`
		GzipLevel int  `koanf:"gzip_level"`
	} `koanf:"server"`

	Params params.Params `koanf:"params"`

	DB core.DBConfig `koanf:"db"`

	View        core.ViewConfig    `koanf:"view"`
	BackendView core.ViewConfig    `koanf:"backendView"`
	Session     core.SessionConfig `koanf:"session"`

	I18N core.I18NConfig `koanf:"i18n"`

	Mail struct {
		Host     string
		Port     int
		Username string
		Password string
	} `koanf:"mail"`

	Captcha core.CaptchaConfig `koanf:"captcha"`

	Migration struct {
		DB     string
		Driver string
		DSN    string
		Path   string
	} `koanf:"migration"`

	CORS struct {
		AllowedOrigins     []string `koanf:"allowed_origins"`
		AllowedHeaders     []string `koanf:"allowed_headers"`
		MaxAge             int      `koanf:"max_age"`
		AllowedCredentials bool     `koanf:"allow_credentials"`
		Debug              bool     `koanf:"debug"`
	} `koanf:"cors"`

	JWT struct {
		SecretKey string `koanf:"secret_key"`
	} `koanf:"jwt"`

	Redis struct {
		Host     string
		Port     int
		Password string
		Database int
	} `koanf:"redis"`
}
