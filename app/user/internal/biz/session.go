package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Session struct {
	Id           string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
	RefreshToken string
}

type SessionRepo interface {
	SetSession(context.Context, Session) error
}

type SessionUsecase struct {
	repo SessionRepo
	log  *log.Helper
}

func NewSessionUsecase(repo SessionRepo, logger log.Logger) *SessionUsecase {
	return &SessionUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *SessionUsecase) SetSession(ctx context.Context, s Session) error {
	return uc.repo.SetSession(ctx, s)
}
