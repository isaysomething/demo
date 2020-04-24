package api

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.IgnoreUnknownKeys(true)
}

func DecodeQueryParams(dst interface{}, ctx *clevergo.Context) (err error) {
	if err = decoder.Decode(dst, ctx.QueryParams()); err != nil {
		return
	}
	if v, ok := dst.(core.Validatable); ok {
		err = v.Validate()
	}
	return
}

type QueryParams struct {
	Sort      string `json:"sort"`
	Direction string `json:"direction"`
}

func (qp QueryParams) Validate() error {
	return validation.ValidateStruct(&qp,
		validation.Field(&qp.Direction, validation.In("asc", "desc")),
	)
}
