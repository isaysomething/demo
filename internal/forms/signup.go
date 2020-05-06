package forms

import (
	"regexp"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/oldmodels"
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/clevergo/demo/pkg/users"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	regUsername = regexp.MustCompile("^[[:alnum:]_-]{5,}$")
)

type AfterSignupEvent struct {
	User *oldmodels.User
}

type Signup struct {
	db             *sqlex.DB
	user           *users.User
	captchaManager *captchas.Manager
	Email          string `json:"email" xml:"email"`
	Username       string `json:"username" xml:"username"`
	Password       string `json:"password" xml:"password"`
	Captcha        string `json:"captcha"`
	CaptchaID      string `json:"captcha_id"`
	onAfterSignup  []func(AfterSignupEvent)
}

func NewSignup(db *sqlex.DB, user *users.User, captchaManager *captchas.Manager) *Signup {
	return &Signup{
		db:             db,
		user:           user,
		captchaManager: captchaManager,
	}
}

func (s *Signup) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Email, validation.Required, is.Email, validation.By(validations.IsUserEmailTaken(s.db))),
		validation.Field(&s.Username,
			validation.Required,
			validation.Match(regUsername),
			validation.By(validations.IsUsernameTaken(s.db)),
		),
		validation.Field(&s.Password, validation.Required),
		validation.Field(&s.CaptchaID, validation.Required),
		validation.Field(&s.Captcha,
			validation.Required,
			validation.By(validations.Captcha(s.captchaManager, s.CaptchaID, true)),
		),
	)
}

func (su *Signup) RegisterOnAfterSignup(f func(AfterSignupEvent)) {
	su.onAfterSignup = append(su.onAfterSignup, f)
}

func (su *Signup) Handle(ctx *clevergo.Context) (*oldmodels.User, error) {
	if err := ctx.Decode(su); err != nil {
		return nil, err
	}

	if err := su.Validate(); err != nil {
		return nil, err
	}
	user, err := oldmodels.CreateUser(su.db, su.Username, su.Email, su.Password)
	if err != nil {
		return nil, err
	}

	su.afterSignup(user)
	return user, nil
}

func (su *Signup) afterSignup(user *oldmodels.User) {
	event := AfterSignupEvent{user}
	for _, f := range su.onAfterSignup {
		f(event)
	}
}
