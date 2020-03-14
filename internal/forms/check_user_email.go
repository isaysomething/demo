package forms

import (
	"github.com/clevergo/demo/internal/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

type CheckUserEmail struct {
	db    *sqlx.DB
	Email string `json:"email"`
}

func NewCheckUserEmail(db *sqlx.DB) *CheckUserEmail {
	return &CheckUserEmail{db: db}
}

func (f *CheckUserEmail) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Email,
			validation.Required,
			is.Email,
			validation.By(validations.IsUserEmailTaken(f.db)),
		),
	)
}
