package user

import (
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/models"
	validation "github.com/go-ozzo/ozzo-validation"
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
		validation.Field(&qp.State, validation.In("active", "inactive")),
	)
}

func (qp QueryParams) StateNumber() int {
	if qp.State == "inactive" {
		return models.UserStateInactive
	}

	return models.UserStateActive
}
