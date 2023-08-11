package biz

import (
	"context"

	v1 "github.com/nei7/ntube/api/user/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.UserServiceErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type UserRepo interface {
	CreateUser(context.Context, *v1.CreateUserRequest) (*v1.User, error)
	GetUserById(context.Context, *v1.GetUserByIdRequest) (*v1.User, error)
	GetUserByEmail(context.Context, *v1.GetUserByEmailRequest) (*v1.GetUserByEmailReply, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, r *v1.CreateUserRequest) (*v1.User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %s", r.Email)
	return uc.repo.CreateUser(ctx, r)
}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, r *v1.GetUserByEmailRequest) (*v1.GetUserByEmailReply, error) {
	uc.log.WithContext(ctx).Infof("GetUserByEmail: %s", r.Email)
	return uc.repo.GetUserByEmail(ctx, r)
}
