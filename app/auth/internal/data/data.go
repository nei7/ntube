package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nei7/ntube/pkg/bootstrap"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEmailVerifyRepo, bootstrap.NewPgxPool)

type Data struct {
	*Queries
	conn *pgx.Conn
	log  *log.Helper
}

func NewData(c *pgx.Conn, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "email_verify/data"))
	d := &Data{
		Queries: New(c),
		conn:    c,
		log:     log,
	}

	cleanup := func() {
		c.Close(context.Background())
		log.Info("closing data resources")
	}

	return d, cleanup, nil
}
