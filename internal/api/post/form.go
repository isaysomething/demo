package post

import (
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/models"
	validation "github.com/go-ozzo/ozzo-validation"
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
		validation.Field(&qp.State, validation.In("published", "draft")),
	)
}

func (qp QueryParams) StateNumber() int {
	if qp.State == "published" {
		return models.PostStatePublished
	}

	return models.PostStateDraft
}
