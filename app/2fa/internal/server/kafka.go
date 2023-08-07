package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/nei7/ntube/api/email/v1"
	"github.com/nei7/ntube/app/2fa/internal/conf"
	"github.com/nei7/ntube/app/2fa/internal/service"

	"github.com/tx7do/kratos-transport/broker"
	"github.com/tx7do/kratos-transport/transport/kafka"
)

func EmailVerifyCreator() broker.Any { return &v1.EmailVerifyRequest{} }

func NewKafkaServer(c *conf.Server, _ log.Logger, svc *service.EmailVerifyService) *kafka.Server {
	ctx := context.Background()

	srv := kafka.NewServer(
		kafka.WithAddress([]string{c.Kafka.Addr}),
		kafka.WithCodec("json"),
	)

	_ = srv.RegisterSubscriber(ctx,
		"2fa.email_verify.ts",
		"email_verify",
		false,
		RegisterEmailVerifyHandler(svc.VerifyEmail),
		EmailVerifyCreator,
	)

	return srv
}

type EmaiVerifyHandler func(ctx context.Context, topic string, headers broker.Headers, msg *v1.EmailVerifyRequest) error

func RegisterEmailVerifyHandler(fnc EmaiVerifyHandler) broker.Handler {
	return func(ctx context.Context, e broker.Event) error {
		var msg *v1.EmailVerifyRequest = nil

		switch t := e.Message().Body.(type) {
		case []byte:
			msg := &v1.EmailVerifyRequest{}
			if err := json.Unmarshal(t, msg); err != nil {
				return err
			}
		case string:
			msg := &v1.EmailVerifyRequest{}
			if err := json.Unmarshal([]byte(t), msg); err != nil {
				return err
			}
		case *v1.EmailVerifyRequest:
			msg = t

		default:
			return fmt.Errorf("unsupported type %T", t)
		}
		return fnc(ctx, e.Topic(), e.Message().Headers, msg)
	}
}
