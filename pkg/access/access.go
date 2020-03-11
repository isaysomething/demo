package access

import (
	"errors"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/users"
)

type Manager struct {
	enforcer     *casbin.Enforcer
	userManager  *users.Manager
	userIDPrefix string
}

func New(enforcer *casbin.Enforcer, userManager *users.Manager) *Manager {
	return &Manager{
		enforcer:     enforcer,
		userManager:  userManager,
		userIDPrefix: "user_",
	}
}

func (m *Manager) Enforce(user, obj, act string) (bool, error) {
	return m.enforcer.Enforce(user, obj, act)
}

func (m *Manager) Middleware(obj, act string) clevergo.MiddlewareFunc {
	return func(ctx *clevergo.Context) error {
		user, err := m.userManager.Get(ctx.Request, ctx.Response)
		if err != nil {
			return err
		}
		if user.IsGuest() {
			return clevergo.StatusError{http.StatusUnauthorized, errors.New("unauthorized")}
		}

		ok, err := m.Enforce(m.getUserID(user), obj, act)
		if err != nil {
			return err
		}
		if !ok {
			return clevergo.StatusError{http.StatusForbidden, errors.New("you have no access to this page")}
		}

		return nil
	}
}

func (m *Manager) getUserID(user *users.User) string {
	return m.userIDPrefix + user.GetIdentity().GetID()
}
