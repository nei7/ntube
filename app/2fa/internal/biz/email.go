package biz

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/nei7/ntube/app/2fa/internal/conf"
)

type EmailSenderUsecase struct {
	password string
	address  string
	host     string
}

func NewEmailSenderUsecase(config *conf.Email) *EmailSenderUsecase {
	return &EmailSenderUsecase{
		host:     config.Host,
		password: config.Password,
		address:  config.Address,
	}
}

func (uc *EmailSenderUsecase) SendEmail(subject string, content []byte, to []string) error {
	auth := smtp.PlainAuth("", uc.address, uc.password, uc.host)
	e := email.NewEmail()
	e.From = uc.address
	e.Subject = subject
	e.To = to
	e.HTML = content
	return e.Send(fmt.Sprintf("%s:%d", uc.host, 587), auth)
}
