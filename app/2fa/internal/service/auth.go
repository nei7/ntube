package service

import (
	"context"

	v1 "github.com/nei7/ntube/api/2fa/v1"
	"github.com/nei7/ntube/app/2fa/internal/biz"
)

type AuthService struct {
	verifyEmailData *biz.AuthUsecase
}

func NewAuthService(ve *biz.AuthUsecase) *AuthService {
	return &AuthService{
		verifyEmailData: ve,
	}
}

func (s *AuthService) VerifyEmail(ctx context.Context, req *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	return s.verifyEmailData.VerifyEmail(ctx, req)
}
