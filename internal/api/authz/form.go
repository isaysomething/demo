package authz

import (
	"github.com/casbin/casbin/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Form struct {
	enforcer    *casbin.Enforcer
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (f *Form) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Permissions, validation.Required),
	)
}
