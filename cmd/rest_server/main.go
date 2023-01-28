package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
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
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/middlewares"
	"github.com/nei7/ntube/internal/opentelemetry"
	"github.com/nei7/ntube/internal/repo"
	"github.com/nei7/ntube/internal/rest"
	"github.com/nei7/ntube/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go:embed static/*
var content embed.FS

func main() {
	var env, addr string

	flag.StringVar(&env, "env", ".env", "Enviroment variables filename")
	flag.StringVar(&addr, "addr", ":3000", "Server address")
	flag.Parse()

	errC, err := run(env, addr)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env, addr string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(env)

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config db.DBConfig
	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	pool, err := db.NewDBConn(config)
	if err != nil {
		return nil, err
	}

	shutdown, err := opentelemetry.InitProviderWithJaegerExporter("rest_server")
	if err != nil {
		return nil, err
	}

	srv, err := newServer(serverConfig{
		addr:        addr,
		DB:          pool,
		jwtKey:      viper.GetString("JWT_KEY"),
		Logger:      logger,
		middlewares: []func(next http.Handler) http.Handler{middlewares.LoggerMiddleware(*logger)},
	})
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			defer shutdown(ctx)
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

	return errChan, nil
}

type serverConfig struct {
	addr        string
	DB          *pgxpool.Pool
	jwtKey      string
	Metrics     http.Handler
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
	tokenManager := service.NewTokenManager(conf.jwtKey)

	rest.NewUserHandler(userService, tokenManager).Register(router)
	router.Handle("/metrics", promhttp.Handler())
	rest.RegisterOpenAPI(router)

	fsys, _ := fs.Sub(content, "static")
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

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
