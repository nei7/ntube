package repo

import (
	"context"

	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/dto"
)

type UserQuery interface {
	Create(ctx context.Context, params dto.CreateUserParams) (db.User, error)
	Find(ctx context.Context, email string) (db.User, error)
}

type UserRepository struct {
	q *db.Queries
}

func NewUserRepo(d db.DBTX) *UserRepository {
	return &UserRepository{
		q: db.New(d),
	}
}

func (r *UserRepository) Create(ctx context.Context, params dto.CreateUserParams) (db.User, error) {
	defer otelSpan(ctx, "User.Create").End()

	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Password: params.Password,
		Email:    params.Email,
	})

	return user, err
}

func (r *UserRepository) Find(ctx context.Context, email string) (db.User, error) {
	defer otelSpan(ctx, "User.Find").End()

	user, err := r.q.GetUser(ctx, email)

	return user, err
}
