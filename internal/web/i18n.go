package web

import (
	"io/ioutil"
	"path"

	"github.com/clevergo/i18n"
	"github.com/gobuffalo/packr/v2"
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

// FileStore is a file store.
type FileStore struct {
	fs      *packr.Box
	decoder i18n.FileDecoder
}

// NewFileStore returns a file store with the given directory and decoder.
func NewFileStore(directory string, decoder i18n.FileDecoder) *FileStore {
	return &FileStore{
		fs:      packr.New(directory, directory),
		decoder: decoder,
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
