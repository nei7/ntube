package repo

import (
	"context"

	"github.com/nei7/ntube/internal/db"
)

type VideQuery interface {
	Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error)
}

type VideoRepo struct {
	q *db.Queries
}

func NewVideRepo(d db.DBTX) *VideoRepo {
	return &VideoRepo{
		q: db.New(d),
	}
}

func (r *VideoRepo) Create(ctx context.Context, params db.CreateVideoParams) (db.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := r.q.CreateVideo(ctx, params)
	return video, err

}
