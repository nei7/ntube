package service

import (
	"context"

	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/repo"
)

type VideoService struct {
	repo repo.VideQuery
}

func (s *VideoService) Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := s.repo.Create(ctx, params)

	return video, err
}
