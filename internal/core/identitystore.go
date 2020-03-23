package core

import (
	"fmt"
	"strconv"

	"github.com/clevergo/auth"
	"github.com/clevergo/demo/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

// IdentityStore is an identity store.
type IdentityStore struct {
	db *sqlx.DB
}

// NewIdentityStore returns an identity store instance.
func NewIdentityStore(db *sqlx.DB) *IdentityStore {
	return &IdentityStore{db: db}
}

// GetIdentity implements IdentityStore.GetIdentity.
func (is *IdentityStore) GetIdentity(id string) (auth.Identity, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ID %q: %s", id, err)
	}
	user, err := models.GetUser(is.db, intID)
	return user, err
}

// GetIdentityByToken implements IdentityStore.GetIdentityByToken.
func (is *IdentityStore) GetIdentityByToken(token, tokenType string) (auth.Identity, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("123456"), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c, _ := t.Claims.(jwt.MapClaims)
	user, err := models.GetUser(is.db, c["id"])
	fmt.Println(user, err)
	if err != nil {
		return nil, err
	}

	return user, nil
}
