package models

import (
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/pkg/sqlex"
	"github.com/clevergo/strutil"
	"golang.org/x/crypto/bcrypt"
)

// User states.
const (
	UserStateDeleted  = 0
	UserStateInactive = 1
	UserStateActive   = 2
)

func init() {
	gob.Register(User{})
}

type User struct {
	ID                 int64            `db:"id" json:"id"`
	Username           string           `db:"username" json:"username"`
	Email              string           `db:"email" json:"email"`
	VerificationToken  sqlex.NullString `db:"verification_token"`
	HashedPassword     string           `db:"hashed_password"`
	PasswordResetToken sqlex.NullString `db:"password_reset_token"`
	State              int              `db:"state" json:"state"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt          sqlex.NullTime   `db:"updated_at" json:"updated_at"`
	DeletedAt          sqlex.NullTime   `db:"deleted_at" json:"deleted_at"`
}

func (u User) GetID() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}

func (u User) IsActive() bool {
	return u.State == UserStateActive
}

func (u User) IsDeleted() bool {
	return u.State == UserStateDeleted
}

func (u User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func (u *User) ValidatePasswordResetToken(duration int64) error {
	return validateToken(u.PasswordResetToken.String, duration)
}

func (u *User) ValidateVerificationToken(duration int64) error {
	return validateToken(u.VerificationToken.String, duration)
}

func (u *User) VerifyEmail(db *sqlex.DB) error {
	_, err := db.NamedExec("UPDATE users SET verification_token=null, state=:state WHERE id=:id", map[string]interface{}{
		"state": UserStateActive,
		"id":    u.ID,
	})

	return err
}

func (u *User) UpdatePassword(db *sqlex.DB, password string) error {
	password, err := GeneratePassword(password)
	if err != nil {
		return err
	}

	_, err = db.NamedExec(
		"UPDATE users SET hashed_password=:password, password_reset_token=null, updated_at=:updated_at WHERE id = :id",
		map[string]interface{}{
			"id":         u.ID,
			"password":   password,
			"updated_at": time.Now(),
		},
	)
	return err
}

func (u *User) GeneratePasswordResetToken(db *sqlex.DB) error {
	token := generateToken(64)
	_, err := db.NamedExec(
		"UPDATE users SET password_reset_token=:token, updated_at=:updated_at WHERE id = :id",
		map[string]interface{}{
			"id":         u.ID,
			"token":      token,
			"updated_at": time.Now(),
		},
	)

	if err != nil {
		return err
	}

	u.PasswordResetToken = sqlex.ToNullString(token)
	return nil
}

func (u *User) GenerateVerificationToken(db *sqlex.DB) error {
	token := generateToken(64)
	_, err := db.NamedExec(
		"UPDATE users SET verification_token=:token, updated_at=:updated_at WHERE id = :id",
		map[string]interface{}{
			"id":         u.ID,
			"token":      token,
			"updated_at": time.Now(),
		},
	)
	if err != nil {
		return err
	}

	u.VerificationToken = sqlex.ToNullString(token)
	return nil
}

func GeneratePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CreateUser(db *sqlex.DB, username, email, password string) (*User, error) {
	hashedPassword, err := GeneratePassword(password)
	if err != nil {
		return nil, err
	}
	verificationToken := generateToken(64)
	now := time.Now()
	query, args, err := squirrel.Insert("users").SetMap(clevergo.Map{
		"username":           username,
		"email":              email,
		"verification_token": verificationToken,
		"hashed_password":    hashedPassword,
		"state":              UserStateInactive,
		"created_at":         now,
		"updated_at":         sqlex.ToNullTime(now),
	}).ToSql()
	if err != nil {
		return nil, err
	}
	res, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetUser(db, id)
}

func GetUser(db *sqlex.DB, id interface{}) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=?", id)
	return user, err
}

func GetUserByUsername(db *sqlex.DB, username string) (*User, error) {
	u := &User{}
	err := db.Get(u, "SELECT * FROM users WHERE username=?", username)
	return u, err
}

func GetUserByEmail(db *sqlex.DB, email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE email=?", email)
	return user, err
}

func GetUserByVerificationToken(db *sqlex.DB, token string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE verification_token=?", token)
	return user, err
}

func GetUserByPasswordResetToken(db *sqlex.DB, token string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE password_reset_token=?", token)
	return user, err
}

func generateToken(length int) string {
	suffix := fmt.Sprintf("_%d", time.Now().Unix())
	return strutil.Random(length-len(suffix)) + suffix
}

func validateToken(token string, duration int64) error {
	if token == "" {
		return errors.New("empty token")
	}
	idx := strings.LastIndex(token, "_")
	if idx == -1 {
		return errors.New("invalid token")
	}

	createdAt, err := strconv.ParseInt(token[idx+1:], 10, 64)
	if err != nil {
		return errors.New("invalid token")
	}

	now := time.Now().Unix()
	if createdAt < now && (createdAt+duration) >= now {
		return nil
	}
	return errors.New("token expired")
}
