package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-chi/chi"
	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/grpc_service"
	"github.com/nei7/gls/internal/middlewares"
	"github.com/nei7/gls/internal/repo"
	"github.com/nei7/gls/internal/rest"
	"github.com/nei7/gls/internal/service"
	"github.com/nei7/gls/pkg/video"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var env, grpc_addr, http_addr string

	flag.StringVar(&env, "env", ".env", "Enviroment variables filename")
	flag.StringVar(&grpc_addr, "grpc_addr", ":3001", "Server address")
	flag.StringVar(&http_addr, "http_addr", ":3002", "Server address")

	flag.Parse()

	errC, err := run(env, grpc_addr, http_addr)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running %s", err)
	}

}

func run(env, grpc_addr, http_addr string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	logger.Sync()

	viper.SetConfigFile(env)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := db.NewDBConfig()
	pool, err := db.NewDBConn(cfg)
	if err != nil {
		return nil, err
	}

	videoRepo := repo.NewVideRepo(pool)

	videoService := service.NewVideoService(videoRepo)

	tokenManager := service.NewTokenManager(viper.GetString("JWT_KEY"))

	videoServer := grpc_service.NewVideoServer(viper.GetString("VIDEO_STORAGE_PATH"), videoService, tokenManager, logger)

	lis, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		return nil, err
	}

	logging := middlewares.LoggerMiddleware(*logger)

	srv, err := newHttpServer(serverConfig{
		addr:        http_addr,
		Logger:      logger,
		middlewares: []func(next http.Handler) http.Handler{logging},
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
			_ = logger.Sync()
			pool.Close()
			defer cancel()
			stop()
			close(errChan)
		}()
		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		if err := lis.Close(); err != nil {
			errChan <- err
		}
		logger.Info("Shutdown completed")

	}()

	go func() {
		logger.Info("Listening and serving", zap.String("grpc_address", grpc_addr), zap.String("http_address", http_addr))

		if err := newGRPCServer(*videoServer, lis); err != nil {
			errChan <- err
		}

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return nil, nil
}

func newGRPCServer(videoServer grpc_service.VideoServer, lis net.Listener) error {
	grpcServer := grpc.NewServer()
	video.RegisterVideoUploadServiceServer(grpcServer, &videoServer)

	return grpcServer.Serve(lis)
}

type serverConfig struct {
	addr        string
	Logger      *zap.Logger
	middlewares []func(next http.Handler) http.Handler
}

func newHttpServer(conf serverConfig) (*http.Server, error) {
	router := chi.NewRouter()

	for _, mw := range conf.middlewares {
		router.Use(mw)
	}

	rest.NewVideoHandler(viper.GetString("VIDEO_STORAGE_PATH")).Register(router)

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
