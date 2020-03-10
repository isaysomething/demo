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

type Store struct {
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

func NewStore(identityStore auth.IdentityStore) *Store {
	return &Store{
		identityStore:    identityStore,
		authIDParam:      "_auth_id",
		authTimeout:      15 * time.Minute,
		authTimeoutParam: "_auth_expire",
	}
}

func (s *Store) SetSessionManager(manager *scs.SessionManager) {
	s.sessionManager = manager
}

func (s *Store) Get(r *http.Request, w http.ResponseWriter) (user *User, err error) {
	user = s.userFromContext(r.Context())
	if user == nil {
		user, err = s.userFromSession(r, w)
		if err != nil {
			log.Println(err)
		}
		if user == nil {
			user = NewUser(s, nil)
		}
		*r = *r.WithContext(context.WithValue(r.Context(), s, user))
	}
	return
}

func (s *Store) userFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(s).(*User)
	return user
}

func (s *Store) userFromSession(r *http.Request, w http.ResponseWriter) (*User, error) {
	if s.sessionManager == nil {
		return nil, errors.New("session was disabled")
	}
	ctx := r.Context()
	if !s.sessionManager.Exists(ctx, s.authIDParam) {
		return nil, errors.New("no auth")
	}
	id := s.sessionManager.GetString(ctx, s.authIDParam)
	if s.authTimeout > 0 {
		expire := s.sessionManager.GetTime(ctx, s.authTimeoutParam)
		if time.Now().After(expire) {
			return nil, errors.New("auth expired")
		}
	}
	if s.sessionManager.Exists(ctx, s.absoluteAuthTimeoutParam) {
		absoluteAuthTimeout := s.sessionManager.GetTime(ctx, s.absoluteAuthTimeoutParam)
		if time.Now().After(absoluteAuthTimeout) {
			return nil, errors.New("auth expired")
		}
	}
	identity, err := s.identityStore.GetIdentity(id)
	if err != nil {
		return nil, err
	}

	if s.authTimeout > 0 {
		s.sessionManager.Put(r.Context(), s.authTimeoutParam, time.Now().Add(s.authTimeout))
	}

	return NewUser(s, identity), nil
}

func (s *Store) RegisterOnBeforeLogin(f func(*LoginEvent) error) {
	s.onBeforeLogin = append(s.onBeforeLogin, f)
}

func (s *Store) RegisterOnAfterLogin(f func(*LoginEvent)) {
	s.onAfterLogin = append(s.onAfterLogin, f)
}

func (s *Store) RegisterOnBeforeLogout(f func(*LogoutEvent) error) {
	s.onBeforeLogout = append(s.onBeforeLogout, f)
}

func (s *Store) RegisterOnAfterLogout(f func(*LogoutEvent)) {
	s.onAfterLogout = append(s.onAfterLogout, f)
}

func (s *Store) Login(r *http.Request, w http.ResponseWriter, user *User, duration time.Duration) (err error) {
	if err = s.beforeLogin(user, duration); err != nil {
		return
	}

	*r = *r.WithContext(context.WithValue(r.Context(), s, user))
	if s.sessionManager != nil {
		ctx := r.Context()
		s.sessionManager.Put(ctx, s.authIDParam, user.identity.GetID())
		if s.authTimeout > 0 {
			s.sessionManager.Put(ctx, s.authTimeoutParam, time.Now().Add(s.authTimeout))
		}
		if duration > 0 {
			s.sessionManager.Put(ctx, s.absoluteAuthTimeoutParam, time.Now().Add(duration))
		}
	}

	s.afterLogin(user, duration)
	return nil
}

func (s *Store) beforeLogin(user *User, duration time.Duration) (err error) {
	event := &LoginEvent{user, duration}
	for _, f := range s.onBeforeLogin {
		if err = f(event); err != nil {
			return
		}
	}

	return
}

func (s *Store) afterLogin(user *User, duration time.Duration) {
	event := &LoginEvent{user, duration}
	for _, f := range s.onAfterLogin {
		f(event)
	}
}

func (s *Store) Logout(r *http.Request, w http.ResponseWriter) (err error) {
	user, err := s.Get(r, w)
	if err != nil {
		return err
	}
	if user.IsGuest() {
		return
	}

	if err = s.beforeLogout(user); err != nil {
		return
	}

	if s.sessionManager != nil {
		s.sessionManager.Remove(r.Context(), s.authIDParam)
		s.sessionManager.Remove(r.Context(), s.authTimeoutParam)
	}

	s.afterLogout(user)
	return nil
}

func (s *Store) beforeLogout(user *User) (err error) {
	event := &LogoutEvent{user}
	for _, f := range s.onBeforeLogout {
		if err = f(event); err != nil {
			return
		}
	}

	return
}

func (s *Store) afterLogout(user *User) {
	event := &LogoutEvent{user}
	for _, f := range s.onAfterLogout {
		f(event)
	}
}
