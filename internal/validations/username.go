package validations

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/clevergo/demo/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var ErrUsername = errors.New("username must contain at least three characters and digits")

var usernameRegexp = regexp.MustCompile(`^[[:alnum:]]{3,}$`)

func ValidateUsername(value interface{}) error {
	s := value.(string)
	if !usernameRegexp.MatchString(s) {
		return ErrUsername
	}
	return nil
}

func IsUsernameAvailable(value interface{}) error {
	s := value.(string)
	exists, err := models.Users(qm.Where("username=?", s)).Exists(context.TODO(), boil.GetContextDB())
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("username %q has been taken", s)
	}
	return nil
}

func IsEmailAvailable(value interface{}) error {
	s := value.(string)
	exists, err := models.Users(qm.Where("email=?", s)).Exists(context.TODO(), boil.GetContextDB())
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email %q has been taken", s)
	}
	return nil
}
