package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/middlewares"
	"github.com/nei7/gls/internal/repo"
	"github.com/nei7/gls/internal/rest"
	"github.com/nei7/gls/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	var env, addr string

	flag.StringVar(&env, "env", ".env", "Enviroment variables filename")
	flag.StringVar(&addr, "addr", ":3000", "Server address")
	flag.Parse()

	err, errC := run(env, addr)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env, addr string) (error, <-chan error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return err, nil
	}

	viper.SetConfigFile(env)

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return err, nil
	}

	var config db.DBConfig
	if err = viper.Unmarshal(&config); err != nil {
		return err, nil
	}

	pool, err := db.NewDBConn(config)
	if err != nil {
		return err, nil
	}

	logging := middlewares.LoggerMiddleware(*logger)

	srv, err := newServer(serverConfig{
		addr:        addr,
		DB:          pool,
		Logger:      logger,
		middlewares: []func(next http.Handler) http.Handler{logging},
	})
	if err != nil {
		return err, nil
	}

	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			_ = logger.Sync()
			pool.Close()
			stop()
			close(errChan)
			cancel()
		}()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}
		logger.Info("Shutdown completed")

	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}

	}()

	return nil, errChan
}

type serverConfig struct {
	addr        string
	DB          *pgxpool.Pool
	Logger      *zap.Logger
	middlewares []func(next http.Handler) http.Handler
}

func newServer(conf serverConfig) (*http.Server, error) {
	router := chi.NewRouter()

	for _, mw := range conf.middlewares {
		router.Use(mw)
	}

	userRepo := repo.NewUserRepo(conf.DB)
	userService := service.NewUserService(conf.Logger, userRepo)

	rest.NewUserHandler(userService).Register(router)

	limiter := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	rateLimitHandler := tollbooth.LimitHandler(limiter, router)

	return &http.Server{
		Handler:           rateLimitHandler,
		Addr:              conf.addr,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
