package listeners

import (
	"fmt"
	"log"

	"github.com/clevergo/demo/internal/forms"
	"github.com/go-mail/mail"
)

type SignUp struct {
	mailer *mail.Dialer
}

func NewSignUp(mailer *mail.Dialer) *SignUp {
	return &SignUp{mailer}
}

func (s *SignUp) AfterSignUp(event forms.AfterSignUpEvent) {
	msg := mail.NewMessage()
	msg.SetHeader("From", s.mailer.Username)
	msg.SetHeader("To", event.User.Email)
	msg.SetHeader("Subject", "Please verify your email address")
	link := "http://localhost:8080/verify-email?verification_token=" + event.User.VerificationToken.String
	msg.SetBody("text/html", fmt.Sprintf(`<a href="%s">%s</a>`, link, link))
	if err := s.mailer.DialAndSend(msg); err != nil {
		log.Println(err)
	}
}
