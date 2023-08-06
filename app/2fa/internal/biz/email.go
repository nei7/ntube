package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/nei7/ntube/api/email/v1"
)

type EmailVerifyRepo interface {
	CreateVerifyEmail(context.Context, *v1.EmailVerifyRequest) (*v1.EmailVerify, error)
}

type EmailVerifyUsecase struct {
	repo EmailVerifyRepo
	log  *log.Helper
}

func NewEmailVerifyUsecase(repo EmailVerifyRepo, logger log.Logger) *EmailVerifyUsecase {
	return &EmailVerifyUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *EmailVerifyUsecase) CreateVerifyEmail(ctx context.Context, g *v1.EmailVerifyRequest) (*v1.EmailVerify, error) {
	return uc.repo.CreateVerifyEmail(ctx, g)
}
