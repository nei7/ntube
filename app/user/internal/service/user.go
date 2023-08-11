package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/segmentio/kafka-go"

	email "github.com/nei7/ntube/api/auth/v1"
	v1 "github.com/nei7/ntube/api/user/v1"
	"github.com/nei7/ntube/app/user/internal/biz"
	"github.com/nei7/ntube/app/user/internal/conf"
	"github.com/nei7/ntube/app/user/util"
)

// GreeterService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc *biz.UserUsecase

	kw *kafka.Writer
}

func NewKafkaSender(conf *conf.Server) (*kafka.Writer, error) {
	w := &kafka.Writer{
		Topic:    conf.Kafka.Topic,
		Addr:     kafka.TCP(conf.Kafka.Addr),
		Balancer: &kafka.LeastBytes{},
	}

	return w, nil
}

func NewUserService(uc *biz.UserUsecase, kw *kafka.Writer) *UserService {
	return &UserService{uc: uc, kw: kw}
}

func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.User, error) {
	user, err := s.uc.CreateUser(ctx, in)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return nil, errors.Conflict(v1.UserServiceErrorReason_ALREADY_EXISTS.String(), "Account already exists")
			}
		}
		return nil, err
	}

	if b, err := json.Marshal(email.SendEmailRequest{
		UserId: user.Id,
		Email:  user.Email,
	}); err == nil {
		s.kw.WriteMessages(ctx, kafka.Message{
			Value: b,
		})

	}

	return user, nil
}

func (s *UserService) VerifyPassword(ctx context.Context, r *v1.VerifyPasswordRequest) (*v1.VerifyPasswordReply, error) {
	user, err := s.uc.GetUserByEmail(ctx, &v1.GetUserByEmailRequest{
		Email: r.Email,
	})
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, errors.BadRequest(v1.UserServiceErrorReason_USER_NOT_FOUND.String(), "User doesn't exist")
		}
		return nil, err
	}

	if !util.CheckPasswordHash(r.Password, user.Password) {
		return nil, errors.Unauthorized(v1.UserServiceErrorReason_INVALID_PASSWORD.String(), "Invalid password")
	}

	return nil, nil
}
