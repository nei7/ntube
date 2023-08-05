package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/jackc/pgx/v5/pgconn"

	v1 "github.com/nei7/ntube/api/user/v1"
	"github.com/nei7/ntube/app/user/internal/biz"
)

// GreeterService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc *biz.UserUsecase
}

func NewGreeterService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
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

	return user, nil
}
