package service

import (
	"context"

	v1 "github.com/nei7/ntube/api/email/v1"
	"github.com/nei7/ntube/app/2fa/internal/biz"
	"github.com/tx7do/kratos-transport/broker"
)

type EmailVerifyService struct {
	uc *biz.EmailVerifyUsecase
}

func NewEmailVerifyService(uc *biz.EmailVerifyUsecase) *EmailVerifyService {
	return &EmailVerifyService{uc: uc}
}

func (s *EmailVerifyService) VerifyEmail(ctx context.Context, topic string, headers broker.Headers, msg *v1.EmailVerifyRequest) error {

	_, err := s.uc.CreateVerifyEmail(ctx, msg)

	return err
}
