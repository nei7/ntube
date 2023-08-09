package service

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	v1 "github.com/nei7/ntube/api/2fa/v1"
	"github.com/nei7/ntube/app/2fa/internal/biz"
	"github.com/tx7do/kratos-transport/broker"
)

var temp *template.Template

func init() {
	temp, _ = template.ParseFiles("./email.tmpl")
}

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

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, struct {
		Link       string
		SecretCode string
		ExpiredAt  string
	}{
		Link:       fmt.Sprintf("http://localhost:8001/v1/email/verify?id=%d&secret_code=%s", data.Id, data.SecretCode),
		ExpiredAt:  data.ExpiredAt.String(),
		SecretCode: data.SecretCode,
	})

	return s.emailSender.SendEmail("Verify email", buf.Bytes(), []string{msg.Email})
}
