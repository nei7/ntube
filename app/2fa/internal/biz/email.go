package biz

import (
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/nei7/ntube/app/2fa/internal/conf"
)

type EmailSenderUsecase struct {
	password string
	address  string
	name     string
}

func NewEmailSenderUsecase(config *conf.Email) *EmailSenderUsecase {
	return &EmailSenderUsecase{
		name:     "nei",
		password: config.Password,
		address:  config.Address,
	}
}

func (uc *EmailSenderUsecase) SendEmail(subject string, content string, to []string) error {
	auth := smtp.PlainAuth("", uc.address, uc.password, "smtp.gmail.com")
	e := email.NewEmail()
	e.From = uc.address
	e.Subject = subject
	e.To = to
	e.HTML = []byte(content)
	return e.Send("smtp.gmail.com:587", auth)
}
