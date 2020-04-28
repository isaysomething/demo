package core

type Config struct {
	HTTP      ServerConfig    `koanf:"http"`
	API       ServerConfig    `koanf:"api"`
	Params    Params          `koanf:"params"`
	DB        DBConfig        `koanf:"db"`
	View      ViewConfig      `koanf:"view"`
	Session   SessionConfig   `koanf:"session"`
	I18N      I18NConfig      `koanf:"i18n"`
	Mailer    MailerConfig    `koanf:"mailer"`
	Captcha   CaptchaConfig   `koanf:"captcha"`
	Migration MigrationConfig `koanf:"migration"`
	CORS      CORSConfig      `koanf:"cors"`
	CSRF      CSRFConfig      `koanf:"csrf"`
	JWT       JWTConfig       `koanf:"jwt"`
	Redis     RedisConfig     `koanf:"redis"`
	Log       LogConfig       `koanf:"log"`
}

type ServerConfig struct {
	Addr string `koanf:"addr"`
	Root string `koanf:"root"`
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
	Duration  int    `konaf:"duration"`
}

type I18NConfig struct {
	Path        string `koanf:"path"`
	Fallback    string `koanf:"fallback"`
	Param       string `koanf:"param"`
	CookieParam string `koanf:"cookie_param"`
}
