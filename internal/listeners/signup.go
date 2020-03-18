package listeners

import (
	"fmt"
	"log"

	"github.com/clevergo/demo/internal/forms"
	"github.com/go-mail/mail"
)

func SendVerificationEmail(mailer *mail.Dialer) func(forms.AfterSignupEvent) {
	return func(event forms.AfterSignupEvent) {
		go func() {
			msg := mail.NewMessage()
			msg.SetHeader("From", mailer.Username)
			msg.SetHeader("To", event.User.Email)
			msg.SetHeader("Subject", "Please verify your email address")
			link := "http://localhost:8080/verify-email?verification_token=" + event.User.VerificationToken.String
			msg.SetBody("text/html", fmt.Sprintf(`<a href="%s">%s</a>`, link, link))
			if err := mailer.DialAndSend(msg); err != nil {
				log.Println(err)
			}
		}()
	}
}
