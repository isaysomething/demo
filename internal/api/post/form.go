package post

import (
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/oldmodels"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Form struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	MarkdownContent string `json:"markdown_content"`
	State           int    `json:"state"`
}

func (f *Form) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Title, validation.Required),
		validation.Field(&f.Content, validation.Required),
		validation.Field(&f.MarkdownContent, validation.Required),
		validation.Field(&f.State, validation.In(oldmodels.PostStateDraft, oldmodels.PostStatePublished)),
	)
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
