package forms

import (
	"errors"
	"fmt"
	"log"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/go-mail/mail"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RequestResetPassword struct {
	db             *sqlex.DB
	mailer         *mail.Dialer
	captchaManager *captchas.Manager
	user           *oldmodels.User
	Email          string `json:"email"`
	Captcha        string `json:"captcha"`
	CaptchaID      string `json:"captcha_id"`
}

func NewRequestResetPassword(db *sqlex.DB, mailer *mail.Dialer, captchaManager *captchas.Manager) *RequestResetPassword {
	return &RequestResetPassword{
		db:             db,
		mailer:         mailer,
		captchaManager: captchaManager,
	}
}

func (f *RequestResetPassword) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.CaptchaID, validation.Required),
		validation.Field(&f.Captcha, validation.Required, validation.By(validations.Captcha(f.captchaManager, f.CaptchaID, true))),
		validation.Field(&f.Email,
			validation.Required,
			is.Email,
			validation.By(f.validateEmail),
		),
	)
}

func (f *RequestResetPassword) validateEmail(value interface{}) error {
	user, err := f.getUser()
	if err != nil || user == nil {
		return errors.New("account does not exist")
	}

	return nil
}

func (f *RequestResetPassword) getUser() (*oldmodels.User, error) {
	if f.user == nil {
		user, err := oldmodels.GetUserByEmail(f.db, f.Email)
		if err != nil {
			return nil, err
		}

		f.user = user
	}

	return f.user, nil
}

func (f *RequestResetPassword) Handle(ctx *clevergo.Context) (err error) {
	if err = ctx.Decode(f); err != nil {
		return
	}
	if err = f.Validate(); err != nil {
		return
	}

	user, _ := f.getUser()
	if err = user.GeneratePasswordResetToken(f.db); err != nil {
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", f.mailer.Username)
	msg.SetHeader("To", f.Email)
	msg.SetHeader("Subject", "Reset Password")
	link := "http://localhost:8080/user/reset-password?token=" + user.PasswordResetToken.String
	msg.SetBody("text/html", fmt.Sprintf(`<a href="%s">%s</a>`, link, link))
	if err := f.mailer.DialAndSend(msg); err != nil {
		log.Println(err)
	}
	err = f.mailer.DialAndSend(msg)
	return
}
