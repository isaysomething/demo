package core

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey []byte
	duration  time.Duration
}

func NewJWTManager(cfg JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey: []byte(cfg.SecretKey),
		duration:  time.Duration(cfg.Duration) * time.Second,
	}
}

func (m *JWTManager) New(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(m.duration).Unix(),
		Subject:   userID,
	})

	return token.SignedString(m.secretKey)
}

func (m *JWTManager) Parse(s string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(s, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
