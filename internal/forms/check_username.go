package forms

import (
	"github.com/clevergo/demo/internal/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type CheckUsername struct {
	db       *sqlx.DB
	Username string `json:"username"`
}

func NewCheckUsername(db *sqlx.DB) *CheckUsername {
	return &CheckUsername{db: db}
}

func (f *CheckUsername) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Username,
			validation.Required,
			validation.Length(5, 0),
			validation.By(validations.IsUsernameTaken(f.db)),
		),
	)
}