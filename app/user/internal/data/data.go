package data

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/nei7/ntube/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

type Data struct {
	q   *Queries
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

func NewData(q *pgx.Conn, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "user/data"))
	d := &Data{
		q:   New(q),
		log: log,
	}

	cleanup := func() {
		q.Close(context.Background())
		log.Info("closing data resources")
	}

	return d, cleanup, nil
}
