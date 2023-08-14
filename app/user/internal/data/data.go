package data

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/nei7/ntube/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewPgxPool, NewSessionRepo, NewRedisClient)

type Data struct {
	q   *Queries
	rdb redis.Cmdable
	log *log.Helper
}

func NewPgxPool(config *conf.Data_Database) (*pgx.Conn, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Username, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.Name,
	}

	return pgx.Connect(context.Background(), dsn.String())
}

func NewRedisClient(config *conf.Data_Redis) redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		DB:       int(config.Db),
		Password: config.Password,
		Username: config.Username,
	})
	if client == nil {
		panic("failed opening connection to redis")
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("cant connect to redis")
	}

	return client
}

func NewData(q *pgx.Conn, rdb redis.Cmdable, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "user/data"))
	d := &Data{
		q:   New(q),
		log: log,
		rdb: rdb,
	}

	cleanup := func() {
		q.Close(context.Background())
		log.Info("closing data resources")

	}

	return d, cleanup, nil
}
