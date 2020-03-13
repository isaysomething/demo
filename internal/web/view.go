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
	Debug bool `koanf:"debug"`
}

// ViewManager wraps jet.Set.
type ViewManager struct {
	*jet.Set
	suffix string
}

// GetTemplate appends suffix to the view filename. 
func (m *ViewManager) GetTemplate(name string) (*jet.Template, error) {
	return m.Set.GetTemplate(name + m.suffix)
}

// NewViewManager returns a view manager associative with the given config.
func NewViewManager(box *packr.Box, cfg ViewConfig) *ViewManager {
	m := &ViewManager{
		Set:    jet.NewHTMLSetLoader(packrloader.New(box)),
		suffix: cfg.Suffix,
	}
	m.Set.SetDevelopmentMode(cfg.Debug)
	if len(cfg.Delims) == 2 {
		m.Set.Delims(cfg.Delims[0], cfg.Delims[1])
	}
	return m
}
