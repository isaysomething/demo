package web

type LogConfig struct {
	File     string `koanf:"file"`
	FileMode uint32 `koanf:"file_mode"`
}

type DBConfig struct {
	Driver string `koanf:"driver"`
	DSN    string `koanf:"dsn"`
}

type I18NConfig struct {
	Directory   string `koanf:"directory"`
	Fallback    string `koanf:"fallback"`
	Param       string `koanf:"param"`
	CookieParam string `koanf:"cookie_param"`
}
