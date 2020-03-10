package users

import (
	"net/http"
	"time"

	"github.com/clevergo/auth"
)

type User struct {
	manager  *Manager
	identity auth.Identity
}

func newUser(manager *Manager, identity auth.Identity) *User {
	return &User{manager, identity}
}

func (u *User) GetIdentity() auth.Identity {
	return u.identity
}

func (u *User) IsGuest() bool {
	return u.identity == nil
}

func (u *User) Login(r *http.Request, w http.ResponseWriter, identity auth.Identity, duration time.Duration) error {
	u.identity = identity
	return u.manager.Login(r, w, u, duration)
}

func (u *User) Logout(r *http.Request, w http.ResponseWriter) error {
	return u.manager.Logout(r, w)
}
