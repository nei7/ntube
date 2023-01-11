package service

import (
	"context"

	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/dto"
	"github.com/nei7/gls/internal/repo"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const otelName = "github.com/nei7/gls/internal/service"

type UserService interface {
	Create(ctx context.Context, params dto.CreateUserParams) (db.User, error)
	Find(ctx context.Context, email string) (db.User, error)
}

type userService struct {
	repo repo.UserQuery
}

func NewUserService(logger *zap.Logger, repo repo.UserQuery) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, params dto.CreateUserParams) (db.User, error) {
	defer otelSpan(ctx, "User.Create").End()

	user, err := s.repo.Create(ctx, params)

	return user, err
}

func (s *userService) Find(ctx context.Context, email string) (db.User, error) {
	defer otelSpan(ctx, "User.Find").End()

	user, err := s.repo.Find(ctx, email)

	return user, err
}

func otelSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)

	return span
}
