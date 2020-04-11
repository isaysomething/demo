package api

import (
	"log"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/jsend"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Handle(ctx *clevergo.Context, err error) {
	if err = jsend.Error(ctx.Response, err.Error()); err != nil {
		log.Println(err)
	}
}
