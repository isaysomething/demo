package core

import "github.com/go-mail/mail"

// MailerConfig mailer config.
type MailerConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

// NewMailer returns a mailer.
func NewMailer(cfg MailerConfig) *mail.Dialer {
	mailer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return mailer
}
