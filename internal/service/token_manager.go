package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager interface {
	NewJWT(userID string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type tokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) *tokenManager {
	return &tokenManager{
		signingKey,
	}
}

func (t *tokenManager) NewJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(t.signingKey))
}

func (t *tokenManager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(t.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't get claims from token")
	}

	id, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid userId")
	}

	return id, nil
}

func (t *tokenManager) NewRefreshToken() (string, error) {
	return "", nil
}
