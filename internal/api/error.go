package api

import (
	"log"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/jsend"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Handle(ctx *clevergo.Context, err error) {
	log.Println(err)
	status := http.StatusOK
	if e, ok := err.(clevergo.Error); ok {
		status = e.Status()
	}
	if err = ctx.JSON(status, jsend.NewError(err.Error(), 0, nil)); err != nil {
		log.Println(err)
	}
}
