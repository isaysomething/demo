package forms

import (
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/db"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CheckUsername struct {
	db       *db.DB
	Username string `json:"username"`
}

func NewCheckUsername(db *db.DB) *CheckUsername {
	return &CheckUsername{db: db}
}

func (f *CheckUsername) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Username,
			validation.Required,
			validation.Length(5, 0),
			validation.Match(regUsername),
			validation.By(validations.IsUsernameTaken(f.db)),
		),
	)
}
