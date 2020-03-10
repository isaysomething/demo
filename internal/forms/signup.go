package forms

import (
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
	Captcha        string `json:"captcha" xml:"captcha"`
	CaptchaID      string `json:"captcha_id" xml:"captcha_id"`
	onAfterSignUp  []func(AfterSignUpEvent)
}

func NewSignUp(db *sqlx.DB, user *users.User, captchaManager *captchas.Manager) *SignUp {
	return &SignUp{
		db:             db,
		user:           user,
		captchaManager: captchaManager,
	}
}

func (su *SignUp) Validate() error {
	return validation.ValidateStruct(su,
		validation.Field(&su.Email, validation.Required, is.Email, validation.By(validations.UniqueUserEmail(su.db))),
		validation.Field(&su.Username, validation.Required, validation.By(validations.UniqueUsername(su.db))),
		validation.Field(&su.Password, validation.Required),
		validation.Field(&su.Captcha, validation.Required, validation.By(validations.Captcha(su.captchaManager, su.CaptchaID))),
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
