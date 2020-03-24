package forms

import (
	"errors"
	"time"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/form"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type ResetPassword struct {
	db       *sqlx.DB
	user     *models.User
	Token    string `json:"token"`
	Password string `json:"password"`
}

func NewResetPassword(db *sqlx.DB) *ResetPassword {
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

	return nil
}

func (f *ResetPassword) getUser() (*models.User, error) {
	if f.user == nil {
		user, err := models.GetUserByPasswordResetToken(f.db, f.Token)
		if err != nil {
			return nil, err
		}

		f.user = user
	}

	return f.user, nil
}

func (f *ResetPassword) Handle(ctx *clevergo.Context) (err error) {
	if err = form.Decode(ctx.Request, f); err != nil {
		return
	}
	if err = f.Validate(); err != nil {
		return
	}

	user, _ := f.getUser()
	password, err := models.GeneratePassword(f.Password)
	if err != nil {
		return err
	}
	_, err = f.db.NamedExec(
		"UPDATE users SET hashed_password=:password, password_reset_token=null, updated_at=:updated_at WHERE id = :id",
		map[string]interface{}{
			"id":         user.ID,
			"password":   password,
			"updated_at": time.Now(),
		},
	)

	return
}
