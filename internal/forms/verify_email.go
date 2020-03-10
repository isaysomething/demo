package forms

import (
	"time"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/users"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type VerifyEmail struct {
	db    *sqlx.DB
	user  *users.User
	Token string `json:"verification_token" xml:"verification_token"`
}

func NewVerifyEmail(db *sqlx.DB, user *users.User) *VerifyEmail {
	return &VerifyEmail{
		db:   db,
		user: user,
	}
}

func (ve *VerifyEmail) Validate() error {
	return validation.ValidateStruct(ve,
		validation.Field(&ve.Token, validation.Required),
	)
}

func (ve *VerifyEmail) Handle(ctx *clevergo.Context) error {
	if err := ve.Validate(); err != nil {
		return err
	}

	identity, err := models.GetUserByVerificationToken(ve.db, ve.Token)
	if err != nil {
		return err
	}
	if err := identity.ValidateVerificationToken(600); err != nil {
		return err
	}

	_, err = ve.db.NamedExec("UPDATE users SET verification_token=:token, status=:status WHERE id=:id", map[string]interface{}{
		"token":  nil,
		"status": models.UserStatusActive,
		"id":     identity.ID,
	})
	if err != nil {
		return err
	}

	return ve.user.Login(ctx.Request, ctx.Response, identity, 24*time.Hour)
}
