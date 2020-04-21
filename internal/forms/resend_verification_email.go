package forms

import (
	"errors"
	"fmt"
	"log"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/db"
	"github.com/go-mail/mail"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ResendVerificationEmail struct {
	db             *db.DB
	mailer         *mail.Dialer
	captchaManager *captchas.Manager
	user           *models.User
	Email          string `valid:"required,email" json:"email" xml:"email"`
	Captcha        string `valid:"required" json:"captcha" xml:"captcha"`
	CaptchaID      string `valid:"required" json:"captcha_id" xml:"captcha_id"`
}

func NewResendVerificationEmail(db *db.DB, mailer *mail.Dialer, captchaManager *captchas.Manager) *ResendVerificationEmail {
	return &ResendVerificationEmail{
		db:             db,
		mailer:         mailer,
		captchaManager: captchaManager,
	}
}

func (f *ResendVerificationEmail) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Captcha, validation.Required, validation.By(validations.Captcha(f.captchaManager, f.CaptchaID, true))),
		validation.Field(&f.Email, validation.Required, is.Email, validation.By(f.validateEmail)),
	)
}

func (f *ResendVerificationEmail) validateEmail(value interface{}) error {
	user, err := f.getUser()
	if err != nil {
		return err
	}

	if user.IsActive() {
		return errors.New("the email has been verified.")
	}

	return nil
}

func (f *ResendVerificationEmail) getUser() (*models.User, error) {
	if f.user == nil {
		user, err := models.GetUserByEmail(f.db, f.Email)
		if err != nil {
			return nil, err
		}
		f.user = user
	}

	return f.user, nil
}

func (f *ResendVerificationEmail) Handle(ctx *clevergo.Context) (err error) {
	if err = ctx.Decode(f); err != nil {
		return
	}
	if err = f.Validate(); err != nil {
		return
	}

	user, _ := f.getUser()
	if err = user.GenerateVerificationToken(f.db); err != nil {
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", f.mailer.Username)
	msg.SetHeader("To", f.Email)
	msg.SetHeader("Subject", "Please verify your email address")
	link := "http://localhost:8080/user/verify-email?token=" + user.VerificationToken.String
	msg.SetBody("text/html", fmt.Sprintf(`<a href="%s">%s</a>`, link, link))
	if err := f.mailer.DialAndSend(msg); err != nil {
		log.Println(err)
	}
	err = f.mailer.DialAndSend(msg)
	return
}
