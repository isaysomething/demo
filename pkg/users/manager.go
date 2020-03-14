package users

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/auth"
)

func init() {
	gob.Register(time.Time{})
}

type Manager struct {
	identityStore            auth.IdentityStore
	sessionManager           *scs.SessionManager
	authIDParam              string
	authTimeout              time.Duration
	authTimeoutParam         string
	absoluteAuthTimeout      time.Duration
	absoluteAuthTimeoutParam string
	onBeforeLogin            []func(*LoginEvent) error
	onAfterLogin             []func(*LoginEvent)
	onBeforeLogout           []func(*LogoutEvent) error
	onAfterLogout            []func(*LogoutEvent)
}

func New(identityStore auth.IdentityStore) *Manager {
	return &Manager{
		identityStore:    identityStore,
		authIDParam:      "_auth_id",
		authTimeout:      15 * time.Minute,
		authTimeoutParam: "_auth_expire",
	}
}

func (m *Manager) SetSessionManager(manager *scs.SessionManager) {
	m.sessionManager = manager
}

func (m *Manager) Get(r *http.Request, w http.ResponseWriter) (user *User, err error) {
	user = m.userFromContext(r.Context())
	if user == nil {
		user, err = m.userFromSession(r, w)
		if err != nil {
			log.Println(err)
		}
		if user == nil {
			user = newUser(m, nil)
		}
		*r = *r.WithContext(context.WithValue(r.Context(), m, user))
	}
	return
}

func (m *Manager) userFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(m).(*User)
	return user
}

func (m *Manager) userFromSession(r *http.Request, w http.ResponseWriter) (*User, error) {
	if m.sessionManager == nil {
		return nil, errors.New("session was disabled")
	}
	ctx := r.Context()
	if !m.sessionManager.Exists(ctx, m.authIDParam) {
		return nil, errors.New("no auth")
	}
	id := m.sessionManager.GetString(ctx, m.authIDParam)
	if m.authTimeout > 0 {
		expire := m.sessionManager.GetTime(ctx, m.authTimeoutParam)
		if time.Now().After(expire) {
			return nil, errors.New("auth expired")
		}
	}
	if m.sessionManager.Exists(ctx, m.absoluteAuthTimeoutParam) {
		absoluteAuthTimeout := m.sessionManager.GetTime(ctx, m.absoluteAuthTimeoutParam)
		if time.Now().After(absoluteAuthTimeout) {
			return nil, errors.New("auth expired")
		}
	}
	identity, err := m.identityStore.GetIdentity(id)
	if err != nil {
		return nil, err
	}

	if m.authTimeout > 0 {
		m.sessionManager.Put(r.Context(), m.authTimeoutParam, time.Now().Add(m.authTimeout))
	}

	return newUser(m, identity), nil
}

func (m *Manager) RegisterOnBeforeLogin(f func(*LoginEvent) error) {
	m.onBeforeLogin = append(m.onBeforeLogin, f)
}

func (m *Manager) RegisterOnAfterLogin(f func(*LoginEvent)) {
	m.onAfterLogin = append(m.onAfterLogin, f)
}

func (m *Manager) RegisterOnBeforeLogout(f func(*LogoutEvent) error) {
	m.onBeforeLogout = append(m.onBeforeLogout, f)
}

func (m *Manager) RegisterOnAfterLogout(f func(*LogoutEvent)) {
	m.onAfterLogout = append(m.onAfterLogout, f)
}

func (m *Manager) Login(r *http.Request, w http.ResponseWriter, user *User, duration time.Duration) (err error) {
	if err = m.beforeLogin(user, duration); err != nil {
		return
	}

	*r = *r.WithContext(context.WithValue(r.Context(), m, user))
	if m.sessionManager != nil {
		ctx := r.Context()
		m.sessionManager.Put(ctx, m.authIDParam, user.identity.GetID())
		if m.authTimeout > 0 {
			m.sessionManager.Put(ctx, m.authTimeoutParam, time.Now().Add(m.authTimeout))
		}
		if duration > 0 {
			m.sessionManager.Put(ctx, m.absoluteAuthTimeoutParam, time.Now().Add(duration))
		}
	}

	m.afterLogin(user, duration)
	return nil
}

func (m *Manager) beforeLogin(user *User, duration time.Duration) (err error) {
	event := &LoginEvent{user, duration}
	for _, f := range m.onBeforeLogin {
		if err = f(event); err != nil {
			return
		}
	}

	return
}

func (m *Manager) afterLogin(user *User, duration time.Duration) {
	event := &LoginEvent{user, duration}
	for _, f := range m.onAfterLogin {
		f(event)
	}
}

func (m *Manager) Logout(r *http.Request, w http.ResponseWriter) (err error) {
	user, err := m.Get(r, w)
	if err != nil {
		return err
	}
	if user.IsGuest() {
		return
	}

	if err = m.beforeLogout(user); err != nil {
		return
	}

	if m.sessionManager != nil {
		m.sessionManager.Remove(r.Context(), m.authIDParam)
		m.sessionManager.Remove(r.Context(), m.authTimeoutParam)
	}

	m.afterLogout(user)
	return nil
}

func (m *Manager) beforeLogout(user *User) (err error) {
	event := &LogoutEvent{user}
	for _, f := range m.onBeforeLogout {
		if err = f(event); err != nil {
			return
		}
	}

	return
}

func (m *Manager) afterLogout(user *User) {
	event := &LogoutEvent{user}
	for _, f := range m.onAfterLogout {
		f(event)
	}
}
