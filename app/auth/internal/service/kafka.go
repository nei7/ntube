package service

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	v1 "github.com/nei7/ntube/api/auth/v1"
	"github.com/nei7/ntube/app/auth/internal/biz"
	"github.com/nei7/ntube/app/auth/internal/conf"
	"github.com/tx7do/kratos-transport/broker"
)

var temp *template.Template

func init() {
	temp, _ = template.ParseFiles("./email.tmpl")
}

type EmailJobService struct {
	authData    *biz.AuthUsecase
	emailSender *biz.EmailSenderUsecase
	url         string
}

func NewEmailJobService(ev *biz.AuthUsecase, es *biz.EmailSenderUsecase, c *conf.Server) *EmailJobService {
	return &EmailJobService{authData: ev, emailSender: es, url: c.Url}
}

func (s *EmailJobService) SendVerifyEmail(ctx context.Context, topic string, headers broker.Headers, msg *v1.SendEmailRequest) error {

	data, err := s.authData.CreateVerifyEmail(ctx, msg)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, struct {
		Id         string
		Link       string
		SecretCode string
		ExpiredAt  string
	}{
		Link:       fmt.Sprintf("email?id=%d&secret_code=%s", data.Id, data.SecretCode),
		ExpiredAt:  data.ExpiredAt.String(),
		SecretCode: data.SecretCode,
	})

	return s.emailSender.SendEmail("Verify email", buf.Bytes(), []string{msg.Email})
}
