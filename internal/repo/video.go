package repo

import (
	"context"

	"github.com/nei7/ntube/internal/datastruct"
	"github.com/nei7/ntube/internal/db"
)

type VideoQuery interface {
	Create(ctx context.Context, params db.CreateVideoParams) (datastruct.Video, error)
}

type VideoRepo struct {
	q *db.Queries
}

func NewVideRepo(d db.DBTX) *VideoRepo {
	return &VideoRepo{
		q: db.New(d),
	}
}

func (r *VideoRepo) Create(ctx context.Context, params db.CreateVideoParams) (datastruct.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := r.q.CreateVideo(ctx, params)
	if err != nil {
		return datastruct.Video{}, err
	}

	user, err := r.q.GetUserById(ctx, params.OwnerID)
	if err != nil {
		return datastruct.Video{}, err
	}

	return datastruct.Video{
		ID:         video.ID.String(),
		Title:      video.Title,
		Thumbnail:  video.Thumbnail,
		Path:       video.Path,
		UploadedAt: video.UploadedAt.Time.Unix(),
		User: datastruct.User{
			Username:   user.Username,
			Created_at: user.CreatedAt.Time.Unix(),
		},
	}, nil

}
