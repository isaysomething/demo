package cmd

import (
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/params"
)

type Config struct {
	Server struct {
		Addr string `koanf:"addr"`
		Root string `koanf:"root"`

		SSL         bool   `koanf:"ssl"`
		SSLCertFile string `koanf:"ssl_cert_file"`
		SSLKeyFile  string `koanf:"ssl_key_file"`

		Log web.LogConfig `koanf:"log"`

		AccessLog         bool   `koanf:"access_log"`
		AccessLogFile     string `koanf:"access_log_file"`
		AccessLogFileMode uint32 `koanf:"access_log_file_mode"`

		Gzip      bool `koanf:"gzip"`
		GzipLevel int  `koanf:"gzip_level"`
	} `koanf:"server"`

	Params params.Params `koanf:"params"`

	DB web.DBConfig `koanf:"db"`

	View        web.ViewConfig    `koanf:"view"`
	BackendView web.ViewConfig    `koanf:"backendView"`
	Session     web.SessionConfig `koanf:"session"`

	I18N web.I18NConfig `koanf:"i18n"`

	Mail struct {
		Host     string
		Port     int
		Username string
		Password string
	} `koanf:"mail"`

	Captcha web.CaptchaConfig `koanf:"captcha"`

	Migration struct {
		DB     string
		Driver string
		DSN    string
		Path   string
	} `koanf:"migration"`

	TencentCaptcha struct {
		SecretID     string `koanf:"secret_id"`
		SecretKey    string `koanf:"secret_key"`
		AppID        uint64 `koanf:"app_id"`
		AppSecretKey string `koanf:"app_secret_key"`
	} `koanf:"tencentCaptcha"`
}
