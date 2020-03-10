package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func NewSession(userID int64, duration int64) Session {
	return Session{
		Token:     GenerateSessionToken(userID),
		UserID:    userID,
		ExpiredAt: time.Now().Add(time.Second * 1),
	}
}

type Session struct {
	Token     string    `db:"token"`
	UserID    int64     `db:"user_id"`
	IPAddress string    `db:"ip_address"`
	UserAgent string    `db:"user_agent"`
	ExpiredAt time.Time `db:"expired_at"`
	CreatedAt time.Time `db:"created_at"`

	user *User
}

func (s *Session) IsExpired() bool {
	return time.Now().Before(s.ExpiredAt)
}

func (s *Session) GetUser(db *sqlx.DB) (*User, error) {
	if s.user == nil {
		user, err := GetUser(db, s.UserID)
		if err != nil {
			return nil, err
		}
		s.user = user
	}

	return s.user, nil
}

func GenerateSessionToken(userID int64) string {
	return ""
}

func GetSession(db *sqlx.DB, token string) (*Session, error) {
	session := &Session{}
	err := db.Get(session, "SELECT * FROM sessions WHERE token=$1", token)
	return session, err
}
