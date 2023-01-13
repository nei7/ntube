package service

import (
	"context"

	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/repo"
)

type VideoService interface {
	Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error)
}

type videoService struct {
	repo repo.VideQuery
}

func NewVideoService(repo repo.VideQuery) VideoService {
	return &videoService{
		repo,
	}
}

func (s *videoService) Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := s.repo.Create(ctx, params)

	return video, err
}
