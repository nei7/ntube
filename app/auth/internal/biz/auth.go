package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/nei7/ntube/api/auth/v1"
)

type AuthRepo interface {
	CreateVerifyEmail(context.Context, *v1.SendEmailRequest) (*v1.EmailVerify, error)
	VerifyEmail(context.Context, *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error)
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) CreateVerifyEmail(ctx context.Context, g *v1.SendEmailRequest) (*v1.EmailVerify, error) {
	return uc.repo.CreateVerifyEmail(ctx, g)
}

func (uc *AuthUsecase) VerifyEmail(ctx context.Context, r *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	return uc.repo.VerifyEmail(ctx, r)
}
