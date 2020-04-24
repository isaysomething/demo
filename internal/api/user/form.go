package user

import (
	"github.com/clevergo/demo/internal/api"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type QueryParams struct {
	api.QueryParams
	Username string `json:"username"`
	Email    string `json:"Email"`
	State    string `json:"state"`
}

func (qp QueryParams) Validate() error {
	return validation.ValidateStruct(&qp,
		validation.Field(&qp.Sort, validation.Required, validation.In("created_at", "updated_at")),
		validation.Field(&qp.State, is.Digit),
	)
}
