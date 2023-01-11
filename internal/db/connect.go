package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
}

func NewDBConfig() DBConfig {
	return DBConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
	}
}

func NewDBConn(config DBConfig) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Username, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.DBName,
	}
	q := dsn.Query()
	q.Add("sslmode", config.SSLMode)

	dsn.RawQuery = q.Encode()

	return pgxpool.Connect(context.Background(), dsn.String())
}
