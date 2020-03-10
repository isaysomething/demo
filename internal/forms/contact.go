package forms

import (
	"html/template"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/form"
	"github.com/go-mail/mail"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Contact struct {
	mailer         *mail.Dialer
	captchaManager *captchas.Manager
	Email          string `json:"email" xml:"email"`
	Subject        string `json:"subject" xml:"subject"`
	Content        string `json:"content" xml:"content"`
	Captcha        string `json:"captcha" xml:"captcha"`
	CaptchaID      string `json:"captcha_id" xml:"captcha_id"`
}

func NewContact(mailer *mail.Dialer, captchaManager *captchas.Manager) *Contact {
	return &Contact{
		mailer:         mailer,
		captchaManager: captchaManager,
	}
}

func (c *Contact) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Captcha, validation.Required, validation.By(validations.Captcha(c.captchaManager, c.CaptchaID))),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Subject, validation.Required),
		validation.Field(&c.Content, validation.Required),
	)
}

func (c *Contact) Handle(ctx *clevergo.Context) (err error) {
	if err = form.Decode(ctx.Request, c); err != nil {
		return
	}
	if err = c.Validate(); err != nil {
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", c.mailer.Username)
	msg.SetHeader("To", c.mailer.Username)
	msg.SetHeader("Subject", c.Subject)
	msg.SetHeader("Reply-To", c.Email)
	msg.SetBody("text/html", template.HTMLEscapeString(c.Content))

	err = c.mailer.DialAndSend(msg)
	return
}
