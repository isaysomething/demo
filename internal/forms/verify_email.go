package forms

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type VerifyEmail struct {
	db    *sqlx.DB
	Token string `json:"token"`
}

func NewVerifyEmail(db *sqlx.DB) *VerifyEmail {
	return &VerifyEmail{
		db: db,
	}
}

func (ve *VerifyEmail) Validate() error {
	return validation.ValidateStruct(ve,
		validation.Field(&ve.Token, validation.Required),
	)
}

func (ve *VerifyEmail) Handle(ctx *clevergo.Context) (err error) {
	if err = ve.Validate(); err != nil {
		return
	}

	identity, err := models.GetUserByVerificationToken(ve.db, ve.Token)
	if err != nil {
		return err
	}
	if err = identity.ValidateVerificationToken(600); err != nil {
		return err
	}

	_, err = ve.db.NamedExec("UPDATE users SET verification_token=null, status=:status WHERE id=:id", map[string]interface{}{
		"status": models.UserStatusActive,
		"id":     identity.ID,
	})
	return
}
