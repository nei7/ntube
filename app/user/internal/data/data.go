package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nei7/ntube/pkg/bootstrap"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, bootstrap.NewPgxPool)

type Data struct {
	q   *Queries
	log *log.Helper
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
