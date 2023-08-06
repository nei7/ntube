package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5/pgtype"

	v1 "github.com/nei7/ntube/api/email/v1"
	"github.com/nei7/ntube/app/2fa/internal/biz"
)

type emailVerifyRepo struct {
	data *Data
	log  *log.Helper
}

func NewEmailVerifyRepo(data *Data, logger log.Logger) biz.EmailVerifyRepo {
	return &emailVerifyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *emailVerifyRepo) CreateVerifyEmail(ctx context.Context, req *v1.EmailVerifyRequest) (*v1.EmailVerify, error) {
	var userId pgtype.UUID

	userId.Scan(req.UserId)

	email, err := r.data.q.CreateVerifyEmail(ctx, CreateVerifyEmailParams{
		Email:      req.Email,
		UserID:     userId,
		SecretCode: "123",
	})
	if err != nil {
		return nil, err
	}

	return &v1.EmailVerify{
		Id: email.ID,
	}, nil
}
