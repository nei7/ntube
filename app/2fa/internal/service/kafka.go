package service

import (
	"context"
	"fmt"

	v1 "github.com/nei7/ntube/api/2fa/v1"
	"github.com/nei7/ntube/app/2fa/internal/biz"
	"github.com/tx7do/kratos-transport/broker"
)

type EmailJobService struct {
	authData    *biz.AuthUsecase
	emailSender *biz.EmailSenderUsecase
}

func NewEmailJobService(ev *biz.AuthUsecase, es *biz.EmailSenderUsecase) *EmailJobService {
	return &EmailJobService{authData: ev, emailSender: es}
}

func (s *EmailJobService) SendVerifyEmail(ctx context.Context, topic string, headers broker.Headers, msg *v1.SendEmailRequest) error {

	data, err := s.authData.CreateVerifyEmail(ctx, msg)
	if err != nil {
		return err
	}

	return s.emailSender.SendEmail("Verify email", fmt.Sprintf("http://localhost:8001/v1/email/verify?id=%d&secret_code=%s", data.Id, data.SecretCode), []string{msg.Email})
}
