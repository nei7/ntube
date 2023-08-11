package data

import (
	"context"
	"fmt"

	"github.com/aidarkhanov/nanoid"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	v1 "github.com/nei7/ntube/api/auth/v1"
	"github.com/nei7/ntube/app/auth/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type emailVerifyRepo struct {
	data *Data
	log  *log.Helper
}

func NewEmailVerifyRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &emailVerifyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *emailVerifyRepo) CreateVerifyEmail(ctx context.Context, req *v1.SendEmailRequest) (*v1.EmailVerify, error) {
	var userId pgtype.UUID

	userId.Scan(req.UserId)

	email, err := r.data.CreateVerifyEmail(ctx, CreateVerifyEmailParams{
		Email:      req.Email,
		UserID:     userId,
		SecretCode: nanoid.New(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.EmailVerify{
		Id:         email.ID,
		ExpiredAt:  timestamppb.New(email.ExpiredAt.Time),
		SecretCode: email.SecretCode,
	}, nil
}

func (r *emailVerifyRepo) VerifyEmail(ctx context.Context, req *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	err := ExecTX(ctx, r.data.conn, func(q *Queries) error {
		result, err := r.data.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         req.Id,
			SecretCode: req.SecretCode,
		})

		if err != nil {
			if errors.Is(pgx.ErrNoRows, err) {
				return errors.BadRequest(v1.AuthServiceErrorReason_EXPIRED_OR_DOESNT_EXISTS.String(), "Link expired or doesn't exists")
			}
			return err
		}

		_, err = r.data.UpdateUser(ctx, UpdateUserParams{
			ID:              result.UserID,
			IsEmailVerified: true,
		})
		if err != nil {
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.VerifyEmailResponse{
		IsVerified: true,
	}, nil
}

func ExecTX(ctx context.Context, conn *pgx.Conn, fn func(q *Queries) error) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(conn)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
