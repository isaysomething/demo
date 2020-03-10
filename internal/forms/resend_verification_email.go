package forms

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/form"
	"github.com/go-mail/mail"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

type ResendVerificationEmail struct {
	db             *sqlx.DB
	mailer         *mail.Dialer
	captchaManager *captchas.Manager
	Email          string `valid:"required,email" json:"email" xml:"email"`
	Captcha        string `valid:"required" json:"captcha" xml:"captcha"`
	CaptchaID      string `valid:"required" json:"captcha_id" xml:"captcha_id"`
}

func NewResendVerificationEmail(db *sqlx.DB, mailer *mail.Dialer, captchaManager *captchas.Manager) *ResendVerificationEmail {
	return &ResendVerificationEmail{
		db:             db,
		mailer:         mailer,
		captchaManager: captchaManager,
	}
}

func (rve *ResendVerificationEmail) Validate() error {
	return validation.ValidateStruct(rve,
		validation.Field(&rve.Captcha, validation.Required, validation.By(validations.Captcha(rve.captchaManager, rve.CaptchaID))),
		validation.Field(&rve.Email, validation.Required, is.Email, validation.By(rve.validateEmail)),
	)
}

func (rve *ResendVerificationEmail) validateEmail(value interface{}) error {
	return nil
}

func (rve *ResendVerificationEmail) Handle(ctx *clevergo.Context) (err error) {
	if err = form.Decode(ctx.Request, rve); err != nil {
		return
	}
	if err = rve.Validate(); err != nil {
		return
	}

	user, err := models.GetUserByEmail(rve.db, rve.Email)
	if err != nil {
		return err
	}
	if user.IsActive() {
		return errors.New("you've been verified your email.")
	}
	newToken := models.GenerateVerificationToken()
	_, err = rve.db.NamedExec("UPDATE users SET verification_token=:token, updated_at=:updated_at", map[string]interface{}{
		"token":      newToken,
		"updated_at": time.Now(),
	})
	if err != nil {
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", rve.mailer.Username)
	msg.SetHeader("To", rve.Email)
	msg.SetHeader("Subject", "Please verify your email address")
	link := "http://localhost:8080/verify-email?verification_token=" + newToken
	msg.SetBody("text/html", fmt.Sprintf(`<a href="%s">%s</a>`, link, link))
	if err := rve.mailer.DialAndSend(msg); err != nil {
		log.Println(err)
	}
	err = rve.mailer.DialAndSend(msg)
	return
}
