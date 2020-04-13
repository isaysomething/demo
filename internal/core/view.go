package core

// ViewData is an alias of map.
type ViewData map[string]interface{}

// ViewConfig contains views manager's settings.
type ViewConfig struct {
	Path   string   `koanf:"path"`
	Suffix string   `koanf:"suffix"`
	Delims []string `koanf:"delims"`
	Debug  bool     `koanf:"debug"`
}
