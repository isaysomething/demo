package web

import (
	"github.com/CloudyKit/jet/v3"
	packrloader "github.com/clevergo/jet-packrloader"
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

type ViewManager struct {
	*jet.Set
	suffix string
}

func (m *ViewManager) GetTemplate(name string) (*jet.Template, error) {
	return m.Set.GetTemplate(name + m.suffix)
}

// NewView returns a views manager with the given config.
func NewViewManager(box *packr.Box, cfg ViewConfig) *ViewManager {
	m := &ViewManager{
		Set:    jet.NewHTMLSetLoader(packrloader.New(box)),
		suffix: cfg.Suffix,
	}
	m.Set.SetDevelopmentMode(true)
	return m
}
