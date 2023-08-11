package bootstrap

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type DBConfig struct {
	Username string
	Password string
	Port     string
	Host     string
	Name     string
}

func NewPgxPool(config *DBConfig) (*pgx.Conn, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Username, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.Name,
	}

	return pgx.Connect(context.Background(), dsn.String())
}
