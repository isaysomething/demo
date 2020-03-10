package cmd

import (
	"github.com/clevergo/demo/internal/web"
	"github.com/clevergo/demo/pkg/params"
)

type Config struct {
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

	Params params.Params `koanf:"params"`

	DB web.DBConfig `koanf:"db"`

	View      web.ViewConfig    `koanf:"view"`
	AdminView web.ViewConfig    `koanf:"adminView"`
	Session   web.SessionConfig `koanf:"session"`

	I18N web.I18NConfig `koanf:"i18n"`

	Mail struct {
		Host     string
		Port     int
		Username string
		Password string
	} `koanf:"mail"`

	Captcha struct {
		Height int
		Width  int
	} `koanf:"captcha"`

	Migration struct {
		DB     string
		Driver string
		DSN    string
		Path   string
	} `koanf:"migration"`
}
