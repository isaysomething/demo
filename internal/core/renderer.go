package core

import (
	"path"
	"reflect"
	"time"

	"github.com/CloudyKit/jet/v3"
	packrloader "github.com/clevergo/jet-packrloader"
	"github.com/clevergo/jetrenderer"
	"github.com/gobuffalo/packr/v2"
)

func NewRenderer(cfg ViewConfig) *jetrenderer.Renderer {
	viewPath := path.Join(cfg.Path, "views")
	box := packr.New(viewPath, viewPath)
	set := jet.NewHTMLSetLoader(packrloader.New(box))
	set.SetDevelopmentMode(cfg.Debug)
	if len(cfg.Delims) == 2 {
		set.Delims(cfg.Delims[0], cfg.Delims[1])
	}

	set.AddGlobalFunc("now", func(_ jet.Arguments) reflect.Value {
		return reflect.ValueOf(time.Now())
	})

	return jetrenderer.New(set)
}
