package post

import (
	"github.com/clevergo/demo/internal/api"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Form struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	State   int    `json:"state"`
}

type QueryParams struct {
	api.QueryParams
	Title string `json:"title"`
	State string `json:"state"`
}

func (qp QueryParams) Validate() error {
	return validation.ValidateStruct(&qp,
		validation.Field(&qp.Sort, validation.Required, validation.In("created_at", "updated_at")),
		validation.Field(&qp.State, is.Digit),
	)
}
