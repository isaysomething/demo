package common

import (
	"fmt"
	"strconv"

	"github.com/clevergo/auth"
	"github.com/clevergo/demo/internal/models"
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
func (is *IdentityStore) GetIdentityByToken(token string) (auth.Identity, error) {
	session, err := models.GetSession(is.db, token)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, fmt.Errorf("token %q is expired", token)
	}

	user, err := session.GetUser(is.db)
	if err != nil {
		return nil, err
	}

	return user, nil
}
