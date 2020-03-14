package validations

import (
	"errors"

	"github.com/clevergo/captchas"
	"github.com/clevergo/demo/pkg/tencentcaptcha"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInvalidCaptcha = errors.New("Invalid captcha")
)

// Captcha validates captcha.
func Captcha(manager *captchas.Manager, id string) validation.RuleFunc {
	return func(value interface{}) error {
		if id == "" {
			return ErrInvalidCaptcha
		}
		captcha, _ := value.(string)
		if err := manager.Verify(id, captcha, true); err != nil {
			return err
		}
		return nil
	}
}

func TencentCaptcha(captcha *tencentcaptcha.Captcha, randstr, ip string) validation.RuleFunc {
	return func(value interface{}) error {
		if ip == "" || randstr == "" {
			return ErrInvalidCaptcha
		}
		ticket, _ := value.(string)
		if _, err := captcha.Validate(ticket, randstr, ip); err != nil {
			return err
		}
		return nil
	}
}
