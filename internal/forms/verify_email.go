package forms

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type VerifyEmail struct {
	db    *sqlx.DB
	user *models.User
	Token string `json:"token"`
}

func NewVerifyEmail(db *sqlx.DB) *VerifyEmail {
	return &VerifyEmail{
		db: db,
	}
}

func (f *VerifyEmail) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Token, validation.Required, validation.By(f.validateUser)),
	)
}

func (f *VerifyEmail) validateUser(value interface{}) error {
	user, err := f.getUser()
	if err != nil {
		return err
	}

	if err = user.ValidateVerificationToken(600); err != nil {
		return err
	}

	return nil
}

func (f *VerifyEmail) getUser() (*models.User, error) {
	if f.user == nil {
		user, err :=  models.GetUserByVerificationToken(f.db, f.Token)
		if err != nil {
			return nil, err
		}

		f.user = user
	}

	return f.user, nil
}

func (f *VerifyEmail) Handle(ctx *clevergo.Context) (err error) {
	if err = f.Validate(); err != nil {
		return
	}

	user, _ := f.getUser()
	return user.VerifyEmail(f.db)
}
