package web

import (
	"path"

	"github.com/clevergo/views/v2"
	"github.com/gobuffalo/packr/v2"
)

// ViewData is an alias of map.
type ViewData map[string]interface{}

// ViewConfig contains views manager's settings.
type ViewConfig struct {
	Path    string   `koanf:"path"`
	Suffix  string   `koanf:"suffix"`
	Delims  []string `koanf:"delims"`
	Layouts []struct {
		Name     string   `koanf:"name"`
		Partials []string `koanf:"partials"`
	} `koanf:"layouts"`
	Cache bool `koanf:"cache"`
}

// NewView returns a views manager with the given config.
func NewView(cfg ViewConfig) *views.Manager {
	viewOpts := []views.Option{
		views.Cache(cfg.Cache),
	}
	if cfg.Suffix != "" {
		viewOpts = append(viewOpts, views.Suffix(cfg.Suffix))
	}
	if len(cfg.Delims) == 2 {
		viewOpts = append(viewOpts, views.Delims(cfg.Delims[0], cfg.Delims[1]))
	}
	viewPath := path.Join(cfg.Path, "views")
	m := views.New(packr.New("views", viewPath), viewOpts...)
	for _, l := range cfg.Layouts {
		m.AddLayout(l.Name, l.Partials...)
	}
	return m
}
