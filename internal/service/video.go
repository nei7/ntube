package service

import (
	"context"

	"github.com/nei7/ntube/internal/datastruct"
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/repo"
)

type VideoService interface {
	Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error)
}

type VideoMessageBrokerRepo interface {
	Created(ctx context.Context, video datastruct.Video) error
	Deleted(ctx context.Context, id string) error
}

type videoService struct {
	repo      repo.VideQuery
	msgBroker VideoMessageBrokerRepo
}

func NewVideoService(repo repo.VideQuery, msgBroker VideoMessageBrokerRepo) VideoService {
	return &videoService{
		repo,
		msgBroker,
	}
}

func (s *videoService) Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := s.repo.Create(ctx, params)

	_ = s.msgBroker.Created(ctx, datastruct.Video{
		ID:          video.ID.String(),
		Path:        video.Path,
		Thumbnail:   video.Thumbnail,
		Description: video.Description,
		Title:       video.Title,
		OwnerID:     video.OwnerID.String(),
		UploadedAt:  video.UploadedAt.Time.Unix(),
	})

	return video, err
}
