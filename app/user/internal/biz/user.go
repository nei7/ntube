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
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) CreateUser(ctx context.Context, g *v1.CreateUserRequest) (*v1.User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %s", g.Email)
	return uc.repo.CreateUser(ctx, g)
}
