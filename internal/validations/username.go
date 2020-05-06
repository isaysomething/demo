package validations

import (
	"errors"
	"regexp"
)

var ErrUsername = errors.New("username must contain at least five characters and digits")

var usernameRegexp = regexp.MustCompile(`^[[:alnum:]]{5,}$`)

func ValidateUsername(value interface{}) error {
	s := value.(string)
	if !usernameRegexp.MatchString(s) {
		return ErrUsername
	}
	return nil
}
