package core

import (
	"fmt"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/auth/authenticators"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/users"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

func NewUserManager(identityStore auth.IdentityStore, sessionManager *scs.SessionManager) *users.Manager {
	m := users.New(identityStore)
	m.SetSessionManager(sessionManager)
	return m
}

// IdentityStore is an identity store.
type IdentityStore struct {
	db *sqlx.DB
}

// NewIdentityStore returns an identity store instance.
func NewIdentityStore(db *sqlx.DB) auth.IdentityStore {
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
		return nil, err
	}
	c, _ := t.Claims.(jwt.MapClaims)
	user, err := models.GetUser(is.db, c["id"])
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthenticator(identityStore auth.IdentityStore) auth.Authenticator {
	return authenticators.NewComposite(
		authenticators.NewBearerToken("api", identityStore),
		authenticators.NewQueryToken("access_token", identityStore),
	)
}
