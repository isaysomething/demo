package forms

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/clevergo/form"
)

type Form interface {
	Handle() (interface{}, error)
}

func Handle(r *http.Request, f Form) (interface{}, error) {
	if err := form.Decode(r, f); err != nil {
		return nil, err
	}

	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}

	return f.Handle()
}
