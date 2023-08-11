package data

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
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

func convertToProto(user User) *v1.User {
	return &v1.User{
		Email:     user.Email,
		Username:  user.Username,
		Id:        uuid.UUIDToString(user.ID),
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
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

	return convertToProto(user), nil
}

func (r *userRepo) GetUserById(ctx context.Context, req *v1.GetUserByIdRequest) (*v1.User, error) {
	var userId pgtype.UUID
	userId.Scan(req.Id)

	user, err := r.data.q.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return convertToProto(user), nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserByEmailReply, error) {
	user, err := r.data.q.GetUser(ctx, req.Email)
	return &v1.GetUserByEmailReply{
		Id:       uuid.UUIDToString(user.ID),
		Password: user.Password,
		Email:    user.Email,
		Username: user.Username,
	}, err
}
