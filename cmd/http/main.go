package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nei7/gls/internal/repo"
	"github.com/nei7/gls/internal/service"
	"go.uber.org/zap"
)

type serverConfig struct {
	addr   string
	DB     *pgxpool.Pool
	Logger *zap.Logger
}

func newServer(conf serverConfig) (*http.Server, error) {
	router := chi.NewRouter()

	userRepo := repo.NewUserRepo(conf.DB)
	userService := service.NewUserService(conf.Logger, userRepo)

}
