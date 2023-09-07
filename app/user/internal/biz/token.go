package biz

import (
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nei7/ntube/app/user/internal/conf"
)

type TokenUsecase struct {
	log       *log.Helper
	secretKey []byte
}

type TokenPayload struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`

	jwt.RegisteredClaims
}

func NewTokenUsecase(c *conf.Token) *TokenUsecase {
	return &TokenUsecase{
		secretKey: []byte(c.Secret),
	}
}

func (uc *TokenUsecase) CreateToken(userId string, sessionId string, duration time.Time) (string, error) {
	claims := TokenPayload{
		SessionId: sessionId,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(duration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.secretKey)
}

func (uc *TokenUsecase) ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &TokenPayload{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(uc.secretKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*TokenPayload)
	if !ok {
		return errors.New("could not parse claims")
	}

	if !claims.ExpiresAt.Before(time.Now()) {
		return errors.New("token is expired")
	}
	return nil
}
