package data

import (
	"context"

	v1 "github.com/nei7/ntube/api/user/v1"
	"github.com/nei7/ntube/app/user/internal/biz"
	"github.com/nei7/ntube/app/user/util"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-kratos/kratos/v2/log"
	uuid "github.com/nei7/ntube/pkg/util"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, g *v1.CreateUserRequest) (*v1.User, error) {

	hashedPassword, err := util.HashPassword(g.Password)
	if err != nil {
		return nil, err
	}

	user, err := r.data.q.CreateUser(ctx, CreateUserParams{
		Username: g.Username,
		Password: hashedPassword,
		Email:    g.Email,
	})

	if err != nil {
		return nil, err
	}

	return &v1.User{
		Email:     user.Email,
		Username:  user.Username,
		Id:        uuid.UUIDToString(user.ID),
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
	}, nil
}
