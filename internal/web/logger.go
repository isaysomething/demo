package web

import (
	"github.com/clevergo/log"
)

type LogConfig struct {
	File     string `koanf:"file"`
	FileMode uint32 `koanf:"file_mode"`
}


type LoggerConfig struct {
	File     string `koanf:"file"`
	FileMode uint32 `koanf:"file_mode"`
}

func NewLogger(cfg LoggerConfig) log.Logger {
	return nil
}