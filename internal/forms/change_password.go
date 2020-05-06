package forms

import (
	"errors"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ChangePassword changes password.
type ChangePassword struct {
	db          *sqlex.DB
	user        *oldmodels.User
	Password    string `json:"password"`     // current password.
	NewPassword string `json:"new_password"` // new password.
}

// NewChangePassword returns a form to change password.
func NewChangePassword(db *sqlex.DB, user *oldmodels.User) *ChangePassword {
	return &ChangePassword{
		db:   db,
		user: user,
	}
}

// Validate validates form data.
func (f *ChangePassword) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.NewPassword, validation.Required),
		validation.Field(&f.Password, validation.Required, validation.By(f.validatePassword)),
	)
}

func (f *ChangePassword) validatePassword(value interface{}) error {
	password, _ := value.(string)
	if err := f.user.ValidatePassword(password); err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

// Handle handles request.
func (f *ChangePassword) Handle(ctx *clevergo.Context) error {
	if err := ctx.Decode(f); err != nil {
		return err
	}

	if err := f.Validate(); err != nil {
		return err
	}

	if err := f.user.UpdatePassword(f.db, f.NewPassword); err != nil {
		return err
	}

	return nil
}
