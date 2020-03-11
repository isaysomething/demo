package web

type LogConfig struct {
	File     string `koanf:"file"`
	FileMode uint32 `koanf:"file_mode"`
}

type I18NConfig struct {
	Path        string `koanf:"path"`
	Fallback    string `koanf:"fallback"`
	Param       string `koanf:"param"`
	CookieParam string `koanf:"cookie_param"`
}
