package users

import (
	"net/http"
	"time"

	"github.com/clevergo/auth"
)

type User struct {
	store    *Store
	identity auth.Identity
}

func NewUser(store *Store, identity auth.Identity) *User {
	return &User{store, identity}
}

func (u *User) GetIdentity() auth.Identity {
	return u.identity
}

func (u *User) IsGuest() bool {
	return u.identity == nil
}

func (u *User) Login(r *http.Request, w http.ResponseWriter, identity auth.Identity, duration time.Duration) error {
	u.identity = identity
	return u.store.Login(r, w, u, duration)
}

func (u *User) Logout(r *http.Request, w http.ResponseWriter) error {
	return u.store.Logout(r, w)
}
