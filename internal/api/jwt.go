package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager() *JWTManager {
	return &JWTManager{}
}

func (m *JWTManager) New(userID int64, duration time.Duration) *jwt.Token {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(time.Second * duration).Unix(),
		Subject:   strconv.FormatInt(userID, 10),
	})
}

func (m *JWTManager) Parse(s string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(s, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.StandardClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return token, nil
}
