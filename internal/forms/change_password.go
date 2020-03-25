package forms

import (
	"time"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/form"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

type ChangePassword struct {
	db          *sqlx.DB
	user        *models.User
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func NewChangePassword(db *sqlx.DB, user *models.User) *ChangePassword {
	return &ChangePassword{
		db:   db,
		user: user,
	}
}

func (cp *ChangePassword) Validate() error {
	return validation.ValidateStruct(cp,
		validation.Field(&cp.NewPassword, validation.Required),
		validation.Field(&cp.Password, validation.Required, validation.By(validations.UserPassword(cp.user))),
	)
}

func (cp *ChangePassword) Handle(ctx *clevergo.Context) error {
	if err := form.Decode(ctx.Request, cp); err != nil {
		return err
	}

	if err := cp.Validate(); err != nil {
		return err
	}

	hashedPassword, err := models.GeneratePassword(cp.NewPassword)
	if err != nil {
		return err
	}

	_, err = cp.db.NamedExec(
		"UPDATE users SET hashed_password=:hashed_password, updated_at=:updated_at WHERE id=:id",
		map[string]interface{}{
			"id":              cp.user.ID,
			"hashed_password": hashedPassword,
			"updated_at":      time.Now(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
