package authz

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Form struct {
	enforcer    *casbin.Enforcer
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}

func (f *Form) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.ID, validation.Required),
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Permissions, validation.Required, validation.Each(validation.By(f.validatePermission))),
	)
}

func (f *Form) validatePermission(value interface{}) error {
	fmt.Println(value)
	return nil
}

type QueryParams struct {
	Name    string   `json:"name"`
	Exclude []string `json:"exclude[]"`
}

func (qp QueryParams) Validate() error {
	return validation.ValidateStruct(&qp)
}
