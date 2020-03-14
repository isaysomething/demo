package web

type I18NConfig struct {
	Path        string `koanf:"path"`
	Fallback    string `koanf:"fallback"`
	Param       string `koanf:"param"`
	CookieParam string `koanf:"cookie_param"`
}
