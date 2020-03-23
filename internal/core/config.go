package core

type Config struct {
	Server    ServerConfig    `koanf:"server"`
	Params    Params          `koanf:"params"`
	DB        DBConfig        `koanf:"db"`
	View      ViewConfig      `koanf:"view"`
	Session   SessionConfig   `koanf:"session"`
	I18N      I18NConfig      `koanf:"i18n"`
	Mailer    MailerConfig    `koanf:"mailer"`
	Captcha   CaptchaConfig   `koanf:"captcha"`
	Migration MigrationConfig `koanf:"migration"`
	CORS      CORSConfig      `koanf:"cors"`
	JWT       JWTConfig       `koanf:"jwt"`
	Redis     RedisConfig     `koanf:"redis"`
	Log       LogConfig       `koanf:"log"`
}

type ServerConfig struct {
	Addr string `koanf:"addr"`
	Root string `koanf:"root"`

	SSL         bool   `koanf:"ssl"`
	SSLCertFile string `koanf:"ssl_cert_file"`
	SSLKeyFile  string `koanf:"ssl_key_file"`

	Log LogConfig `koanf:"log"`

	AccessLog         bool   `koanf:"access_log"`
	AccessLogFile     string `koanf:"access_log_file"`
	AccessLogFileMode uint32 `koanf:"access_log_file_mode"`

	Gzip      bool `koanf:"gzip"`
	GzipLevel int  `koanf:"gzip_level"`
}

type MigrationConfig struct {
	DB     string `koanf:"db"`
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
	Path   string `koanf:"path"`
}

type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	Database int    `koanf:"database"`
}

type JWTConfig struct {
	SecretKey string `koanf:"secret_key"`
}

type I18NConfig struct {
	Path        string `koanf:"path"`
	Fallback    string `koanf:"fallback"`
	Param       string `koanf:"param"`
	CookieParam string `koanf:"cookie_param"`
}
