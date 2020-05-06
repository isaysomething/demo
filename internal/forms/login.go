package forms

import (
	"errors"
	"time"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/clevergo/demo/pkg/users"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Login is a login form.
type Login struct {
	db             *sqlex.DB
	user           *users.User
	identity       *oldmodels.User
	captchaManager *captchas.Manager
	CaptchaID      string `schema:"captcha_id" json:"captcha_id"`
	Captcha        string `schema:"captcha" json:"captcha"`
	Email          string `schema:"email" json:"email"`
	Password       string `schema:"password" json:"password"`
}

// NewLogin returns a login form.
func NewLogin(db *sqlex.DB, user *users.User, captchaManager *captchas.Manager) *Login {
	return &Login{
		db:             db,
		user:           user,
		captchaManager: captchaManager,
	}
}

// Validate validates form data.
func (l *Login) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.CaptchaID, validation.Required),
		validation.Field(&l.Captcha,
			validation.Required,
			validation.By(validations.Captcha(l.captchaManager, l.CaptchaID, true)),
		),
		validation.Field(&l.Password, validation.Required),
		validation.Field(&l.Email,
			validation.Required,
			is.Email,
			validation.By(validation.RuleFunc(l.validateUser)),
		),
	)
}

var (
	errUserIncorrectPassword = errors.New("incorrect username or password")
	errUserInactive          = errors.New("account is not active, please verify your email")
	errUserDeleted           = errors.New("account was deleted")
)

func (l *Login) validateUser(_ interface{}) error {
	identity := l.getIdentity()
	if identity == nil {
		return errUserIncorrectPassword
	}
	if err := identity.ValidatePassword(l.Password); err != nil {
		return errUserIncorrectPassword
	}
	if identity.IsDeleted() {
		return errUserDeleted
	}
	if !identity.IsActive() {
		return errUserInactive
	}
	return nil
}

func (l *Login) getIdentity() *oldmodels.User {
	if l.identity == nil {
		identity, err := oldmodels.GetUserByEmail(l.db, l.Email)
		if err == nil {
			l.identity = identity
		}

	}

	return l.identity
}

// Handle handles login request.
func (l *Login) Handle(ctx *clevergo.Context) (*oldmodels.User, error) {
	if err := ctx.Decode(l); err != nil {
		return nil, err
	}
	identity := l.getIdentity()
	if err := l.user.Login(ctx.Request, ctx.Response, identity, time.Hour); err != nil {
		return nil, err
	}

	return identity, nil
}
