package forms

import (
	"time"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/users"
	"github.com/clevergo/form"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

type Login struct {
	db             *sqlx.DB
	user           *users.User
	identity       *models.User
	captchaManager *captchas.Manager
	CaptchaID      string `json:"captcha_id"`
	Captcha        string `json:"captcha"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ipAddr         string
}

func NewLogin(db *sqlx.DB, user *users.User, captchaManager *captchas.Manager) *Login {
	return &Login{
		db:             db,
		user:           user,
		captchaManager: captchaManager,
	}
}

func (l *Login) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.By(validations.UserPassword(l.getIdentity()))),
		validation.Field(&l.CaptchaID, validation.Required),
		validation.Field(&l.Captcha,
			validation.Required,
			validation.By(validations.Captcha(l.captchaManager, l.CaptchaID)),
		),
	)
}

func (l *Login) ValidatePassword() bool {
	identity := l.getIdentity()
	if identity == nil {
		return false
	}

	err := identity.ValidatePassword(l.Password)
	return err == nil
}

func (l *Login) getIdentity() *models.User {
	if l.identity == nil {
		identity, err := models.GetUserByEmail(l.db, l.Email)
		if err == nil {
			l.identity = identity
		}

	}

	return l.identity
}

func (l *Login) Handle(ctx *clevergo.Context) (*models.User, error) {
	if err := form.Decode(ctx.Request, l); err != nil {
		return nil, err
	}
	/*if _, err := govalidator.ValidateStruct(l); err != nil {
		return nil, err
	}*/
	l.ipAddr = "127.0.0.1"
	if err := l.Validate(); err != nil {
		return nil, err
	}
	identity := l.getIdentity()
	if err := l.user.Login(ctx.Request, ctx.Response, identity, time.Hour); err != nil {
		return nil, err
	}
	return identity, nil
}
