package service

import (
	"context"

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
	return s.uc.CreateUser(ctx, in)
}
