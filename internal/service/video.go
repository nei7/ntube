package service

import (
	"context"
	"time"

	"github.com/mercari/go-circuitbreaker"
	"github.com/nei7/ntube/internal"
	"github.com/nei7/ntube/internal/datastruct"
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/dto"
	"github.com/nei7/ntube/internal/repo"
	"go.uber.org/zap"
)

type VideoService interface {
	Create(ctx context.Context, params db.CreateVideoParams) (datastruct.Video, error)
	Search(ctx context.Context, params dto.VideoSearchParams) (datastruct.SearchResult, error)
}

type TaskSearchRepository interface {
	Search(ctx context.Context, params dto.VideoSearchParams) (datastruct.SearchResult, error)
}

type VideoMessageBrokerRepo interface {
	Created(ctx context.Context, video datastruct.Video) error
	Deleted(ctx context.Context, id string) error
}

type videoService struct {
	repo      repo.VideoQuery
	msgBroker VideoMessageBrokerRepo
	search    TaskSearchRepository
	cb        *circuitbreaker.CircuitBreaker
}

func NewVideoService(logger *zap.Logger, repo repo.VideoQuery, search TaskSearchRepository, msgBroker VideoMessageBrokerRepo) VideoService {
	return &videoService{
		repo:      repo,
		msgBroker: msgBroker,
		search:    search,
		cb: circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(2*time.Minute),
			circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				logger.Info("state changed",
					zap.String("old", string(oldState)),
					zap.String("new", string(newState)),
				)
			}),
		),
	}
}

func (s *videoService) Search(ctx context.Context, params dto.VideoSearchParams) (_ datastruct.SearchResult, err error) {
	defer otelSpan(ctx, "Video.Search").End()

	if !s.cb.Ready() {
		return datastruct.SearchResult{}, internal.NewErrorf(internal.ErrorCodeUnknown, "service not available")
	}

	defer func() {
		s.cb.Done(ctx, err)
	}()

	res, err := s.search.Search(ctx, params)
	if err != nil {
		return datastruct.SearchResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "Search")
	}

	return res, nil
}

func (s *videoService) Create(ctx context.Context, params db.CreateVideoParams) (datastruct.Video, error) {
	defer otelSpan(ctx, "Video.Create").End()

	video, err := s.repo.Create(ctx, params)

	_ = s.msgBroker.Created(ctx, video)

	return video, err
}
