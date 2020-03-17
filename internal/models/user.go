package models

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// User statuses.
const (
	UserStatusDeleted  = 0
	UserStatusInactive = 1
	UserStatusActive   = 10
)

func init() {
	gob.Register(User{})
}

type User struct {
	ID                 int64          `db:"id"`
	Username           string         `db:"username"`
	Email              string         `db:"email"`
	VerificationToken  sql.NullString `db:"verification_token"`
	HashedPassword     string         `db:"hashed_password"`
	PasswordResetToken sql.NullString `db:"password_reset_token"`
	Status             int            `db:"status"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          sql.NullTime   `db:"updated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"`
}

func (u User) GetID() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}

func (u User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func (u User) ValidateVerificationToken(duration int64) error {
	return ValidateVerificationToken(u.VerificationToken.String, duration)
}

func GeneratePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CreateUser(db *sqlx.DB, username, email, password string) (*User, error) {
	hashedPassword, err := GeneratePassword(password)
	if err != nil {
		return nil, err
	}
	verificationToken := GenerateVerificationToken()
	res, err := db.NamedExec(
		`INSERT INTO users (username, email, verification_token, hashed_password, status, created_at) 
		VALUES (:username, :email, :verification_token, :hashed_password, :status, :created_at)`,
		map[string]interface{}{
			"username":           username,
			"email":              email,
			"verification_token": verificationToken,
			"hashed_password":    hashedPassword,
			"status":             UserStatusInactive,
			"created_at":         time.Now(),
		},
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return GetUser(db, id)
}

func GetUser(db *sqlx.DB, id interface{}) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=?", id)
	return user, err
}

func GetUserByUsername(db *sqlx.DB, username string) (*User, error) {
	u := &User{}
	err := db.Get(u, "SELECT * FROM users WHERE username=?", username)
	return u, err
}

func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE email=?", email)
	return user, err
}

func GetUserByVerificationToken(db *sqlx.DB, token string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE verification_token=?", token)
	return user, err
}

func GetUserByPasswordResetToken(db *sqlx.DB, token string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE password_reset_token=?", token)
	return user, err
}

func GenerateVerificationToken() string {
	return fmt.Sprintf("%s_%d", randomString(64), time.Now().Unix())
}

func ValidateVerificationToken(token string, duration int64) error {
	idx := strings.LastIndex(token, "_")
	if idx == -1 {
		return errors.New("invalid token")
	}

	createdAt, err := strconv.ParseInt(token[idx+1:], 10, 64)
	if err != nil {
		return errors.New("invalid token")
	}

	now := time.Now().Unix()
	fmt.Println(now, createdAt)
	if createdAt < now && (createdAt+duration) >= now {
		return nil
	}
	return errors.New("token expired")
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
