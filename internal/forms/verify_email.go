package forms

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type VerifyEmail struct {
	db    *sqlex.DB
	user  *oldmodels.User
	Token string `json:"token"`
}

func NewVerifyEmail(db *sqlex.DB) *VerifyEmail {
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

func (f *VerifyEmail) getUser() (*oldmodels.User, error) {
	if f.user == nil {
		user, err := oldmodels.GetUserByVerificationToken(f.db, f.Token)
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
