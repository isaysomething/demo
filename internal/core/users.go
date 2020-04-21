package core

import (
	"fmt"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
	"github.com/clevergo/auth/authenticators"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/db"
	"github.com/clevergo/demo/pkg/users"
)

func NewUserManager(identityStore auth.IdentityStore, sessionManager *scs.SessionManager) *users.Manager {
	m := users.New(identityStore)
	m.SetSessionManager(sessionManager)
	return m
}

// IdentityStore is an identity store.
type IdentityStore struct {
	db      *db.DB
	manager *JWTManager
}

// NewIdentityStore returns an identity store instance.
func NewIdentityStore(db *db.DB, manager *JWTManager) auth.IdentityStore {
	return &IdentityStore{db: db, manager: manager}
}

// GetIdentity implements IdentityStore.GetIdentity.
func (s *IdentityStore) GetIdentity(id string) (auth.Identity, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ID %q: %s", id, err)
	}
	user, err := models.GetUser(s.db, intID)
	return user, err
}

// GetIdentityByToken implements IdentityStore.GetIdentityByToken.
func (s *IdentityStore) GetIdentityByToken(token, tokenType string) (auth.Identity, error) {
	claims, err := s.manager.Parse(token)
	if err != nil {
		return nil, err
	}
	user, err := models.GetUser(s.db, claims.Subject)
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
