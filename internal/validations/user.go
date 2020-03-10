package validations

import (
	"database/sql"
	"errors"

	"github.com/clevergo/demo/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jmoiron/sqlx"
)

// Errors
var (
	ErrIncorrectPassword = errors.New("incorrect username or password")
	ErrEmailWasTaken     = errors.New("email was taken")
	ErrUsernameWasTaken  = errors.New("username was taken")
)

// UserPassword validates user password.
func UserPassword(user *models.User) validation.RuleFunc {
	return func(value interface{}) error {
		password, _ := value.(string)
		if user == nil {
			return ErrIncorrectPassword
		}
		if err := user.ValidatePassword(password); err != nil {
			return ErrIncorrectPassword
		}
		return nil
	}
}

func UniqueUsername(db *sqlx.DB) validation.RuleFunc {
	return func(value interface{}) error {
		username, _ := value.(string)
		if _, err := models.GetUserByUsername(db, username); err == nil || err != sql.ErrNoRows {
			return ErrUsernameWasTaken
		}
		return nil
	}
}

func UniqueUserEmail(db *sqlx.DB) validation.RuleFunc {
	return func(value interface{}) error {
		email, _ := value.(string)
		if _, err := models.GetUserByEmail(db, email); err == nil || err != sql.ErrNoRows {
			return ErrEmailWasTaken
		}
		return nil
	}
}
