package validations

import (
	"errors"
	"regexp"
)

type passwordRegexp struct {
	*regexp.Regexp
	err error
}

var (
	ErrPasswordDigit              = errors.New("password must contain at least one digit")
	ErrPasswordLowercaseCharacter = errors.New("password must contain at least one lowercase character")
	ErrPasswordUppercaseCharacter = errors.New("password must contain at least one uppercase character")
	ErrPasswordMinLength          = errors.New("password must contain at least six characters and digits")
)

var passwordRegexps = []*passwordRegexp{
	{regexp.MustCompile(`[0-9]+`), ErrPasswordDigit},
	{regexp.MustCompile(`[a-z]+`), ErrPasswordLowercaseCharacter},
	{regexp.MustCompile(`[A-Z]+`), ErrPasswordUppercaseCharacter},
	{regexp.MustCompile(`^[[:alnum:]#@$!]{6,}$`), ErrPasswordMinLength},
}

func ValidatePassword(value interface{}) error {
	password := value.(string)
	for _, reg := range passwordRegexps {
		if !reg.MatchString(password) {
			return reg.err
		}
	}

	return nil
}
