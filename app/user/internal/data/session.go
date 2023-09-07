package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/nei7/ntube/app/user/internal/biz"
)

var _ biz.SessionRepo = (*sessionRepo)(nil)

type sessionRepo struct {
	data *Data
	log  *log.Helper
}

func NewSessionRepo(data *Data, logger log.Logger) biz.SessionRepo {
	return &sessionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *sessionRepo) SetSession(ctx context.Context, s biz.Session) error {
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = r.data.rdb.Set(ctx, s.Id, string(buf), time.Hour).Err()
	if err != nil {
		r.log.Errorf("failed to create a session cache:redis.Set(%v) error(%v)", s, err)

	}

	return err
}

func (r *sessionRepo) GetSession(ctx context.Context, id string) (*biz.Session, error) {
	data, err := r.data.rdb.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NotFound("SESSION_NOT_FOUND", "session not found")
		}

		return nil, err
	}

	var session = &biz.Session{}
	err = json.Unmarshal([]byte(data), session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
