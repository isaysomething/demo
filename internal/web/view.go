package web

import (
	"path"

	"github.com/clevergo/views/v2"
)

type ViewConfig struct {
	Directory string   `koanf:"directory"`
	Suffix    string   `koanf:"suffix"`
	Delims    []string `koanf:"delims"`
	Layouts   []struct {
		Name     string   `koanf:"name"`
		Partials []string `koanf:"partials"`
	} `koanf:"layouts"`
	Cache bool `koanf:"cache"`
}

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
	m := views.New(path.Join(cfg.Directory, "views"), viewOpts...)
	for _, l := range cfg.Layouts {
		m.AddLayout(l.Name, l.Partials...)
	}
	return m
}
