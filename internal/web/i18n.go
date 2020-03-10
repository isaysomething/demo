package web

import (
	"github.com/clevergo/i18n"
)

type I18N struct {
	cfg I18NConfig
	*i18n.Translators
}

func NewI18N(cfg I18NConfig) *I18N {
	return &I18N{
		cfg: cfg,
	}
}
