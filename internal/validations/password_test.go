package validations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		password  string
		shouldErr bool
		err       error
	}{
		{"12Foo", true, ErrPasswordMinLength},
		{"123456", true, ErrPasswordLowercaseCharacter},
		{"123foo", true, ErrPasswordUppercaseCharacter},
		{"FooBar", true, ErrPasswordDigit},
		{"123Foo", false, nil},
		{"Foo#Bar@123", false, nil},
		{"P@$$w0rd!", false, nil},
	}

	for _, c := range cases {
		err := ValidatePassword(c.password)
		if c.shouldErr {
			assert.Equal(t, c.err, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
