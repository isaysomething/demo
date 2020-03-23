package core

import "encoding/gob"

func init() {
	gob.Register(Flashes{})
}

type Flashes []Flash

type Flash interface {
	Message() string
}
