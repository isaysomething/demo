package forms

import (
	"regexp"

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

var (
	regUsername = regexp.MustCompile("^[[:alnum:]_-]{5,}$")
)

type AfterSignUpEvent struct {
	User *models.User
}

type SignUp struct {
	db             *sqlx.DB
	user           *users.User
	captchaManager *captchas.Manager
	Email          string `json:"email" xml:"email"`
	Username       string `json:"username" xml:"username"`
	Password       string `json:"password" xml:"password"`
	Captcha        string `json:"captcha"`
	CaptchaID      string `json:"captcha_id"`
	onAfterSignUp  []func(AfterSignUpEvent)
}

func NewSignUp(db *sqlx.DB, user *users.User, captchaManager *captchas.Manager) *SignUp {
	return &SignUp{
		db:             db,
		user:           user,
		captchaManager: captchaManager,
	}
}

func (s *SignUp) Validate() error {
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
			validation.By(validations.Captcha(s.captchaManager, s.CaptchaID, false)),
		),
	)
}

func (su *SignUp) RegisterOnAfterSignUp(f func(AfterSignUpEvent)) {
	su.onAfterSignUp = append(su.onAfterSignUp, f)
}

func (su *SignUp) Handle(ctx *clevergo.Context) (*models.User, error) {
	if err := form.Decode(ctx.Request, su); err != nil {
		return nil, err
	}

	if err := su.Validate(); err != nil {
		return nil, err
	}
	user, err := models.CreateUser(su.db, su.Username, su.Email, su.Password)
	if err != nil {
		return nil, err
	}

	su.afterSignUp(user)
	return user, nil
}

func (su *SignUp) afterSignUp(user *models.User) {
	event := AfterSignUpEvent{user}
	for _, f := range su.onAfterSignUp {
		f(event)
	}
}
