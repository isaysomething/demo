package validations

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Errors
var (
	ErrIncorrectPassword = errors.New("incorrect username or password")
	ErrEmailWasTaken     = errors.New("email was taken")
	ErrUsernameWasTaken  = errors.New("username was taken")
)

// UserPassword validates user password.
func UserPassword(user *oldmodels.User) validation.RuleFunc {
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

// IsUsernameTaken validates whether the username was taken.
func IsUsernameTaken(db *sqlex.DB) validation.RuleFunc {
	return func(value interface{}) error {
		username, _ := value.(string)
		_, err := oldmodels.GetUserByUsername(db, username)
		if err == nil {
			return fmt.Errorf("username %s was taken", username)
		}
		return nil
	}
}

// IsUserEmailTaken validates whether the user email was taken.
func IsUserEmailTaken(db *sqlex.DB) validation.RuleFunc {
	return func(value interface{}) error {
		email, _ := value.(string)
		if _, err := oldmodels.GetUserByEmail(db, email); err == nil || err != sql.ErrNoRows {
			return ErrEmailWasTaken
		}
		return nil
	}
}
