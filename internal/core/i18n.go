package core

import (
	"io/ioutil"
	"path"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/i18n"
	"github.com/gobuffalo/packr/v2"
)

func NewI18N(cfg I18NConfig, store i18n.Store) (*i18n.Translators, error) {
	i18nOpts := []i18n.Option{i18n.Fallback(cfg.Fallback)}
	translators := i18n.New(i18nOpts...)
	if err := translators.Import(store); err != nil {
		return nil, err
	}

	return translators, nil
}

// FileStore is a file store.
type FileStore struct {
	fs      *packr.Box
	decoder i18n.FileDecoder
}

// NewFileStore returns a file store with the given directory and decoder.
func NewFileStore(cfg I18NConfig) i18n.Store {
	directory := cfg.Path
	return &FileStore{
		fs:      packr.New(directory, directory),
		decoder: &i18n.JSONFileDecoder{},
	}
}

// Get implements Store.Get.
func (s *FileStore) Get() (i18n.Translations, error) {
	translations := make(i18n.Translations)
	for _, name := range s.fs.List() {
		dir := path.Dir(name)

		f, err := s.fs.Open(name)
		if err != nil {
			return translations, err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return translations, err
		}
		var v i18n.Translation
		if err = s.decoder.Decode(data, &v); err != nil {
			return translations, err
		}
		translations[dir] = v
	}

	return translations, nil
}

func NewI18NLanguageParsers(cfg I18NConfig) (parsers []i18n.LanguageParser) {
	if cfg.Param != "" {
		parsers = append(parsers, i18n.NewURLLanguageParser(cfg.Param))
	}
	if cfg.CookieParam != "" {
		parsers = append(parsers, i18n.NewCookieLanguageParser(cfg.CookieParam))
	}
	parsers = append(parsers, i18n.HeaderLanguageParser{})
	return
}

type I18NMiddleware clevergo.MiddlewareFunc

func NewI18NMiddleware(translators *i18n.Translators, parsers []i18n.LanguageParser) I18NMiddleware {
	return I18NMiddleware(clevergo.WrapHH(i18n.Middleware(translators, parsers...)))
}
