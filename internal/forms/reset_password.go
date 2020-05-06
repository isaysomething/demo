package forms

import (
	"errors"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ResetPassword struct {
	db       *sqlex.DB
	user     *oldmodels.User
	Token    string `json:"token"`
	Password string `json:"password"`
}

func NewResetPassword(db *sqlex.DB) *ResetPassword {
	return &ResetPassword{
		db: db,
	}
}

func (f *ResetPassword) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Password,
			validation.Required,
			validation.Length(6, 0),
		),
		validation.Field(&f.Token,
			validation.Required,
			validation.By(f.validateToken),
		),
	)
}

func (f *ResetPassword) validateToken(value interface{}) error {
	user, err := f.getUser()
	if err != nil || user == nil {
		return errors.New("account does not exist")
	}
	if err = user.ValidatePasswordResetToken(600); err != nil {
		return err
	}

	return nil
}

func (f *ResetPassword) getUser() (*oldmodels.User, error) {
	if f.user == nil {
		user, err := oldmodels.GetUserByPasswordResetToken(f.db, f.Token)
		if err != nil {
			return nil, err
		}

		f.user = user
	}

	return f.user, nil
}

func (f *ResetPassword) Handle(ctx *clevergo.Context) (err error) {
	if err = ctx.Decode(f); err != nil {
		return
	}
	if err = f.Validate(); err != nil {
		return
	}

	user, _ := f.getUser()
	return user.UpdatePassword(f.db, f.Password)
}
