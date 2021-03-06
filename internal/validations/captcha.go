package validations

import (
	"errors"

	"github.com/clevergo/captchas"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidCaptcha = errors.New("Invalid captcha")
)

// Captcha validates captcha.
func Captcha(manager *captchas.Manager, id string, clear bool) validation.RuleFunc {
	return func(value interface{}) error {
		if id == "" {
			return ErrInvalidCaptcha
		}
		captcha, _ := value.(string)
		if err := manager.Verify(id, captcha, clear); err != nil {
			return err
		}
		return nil
	}
}
