package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/grpc_service"
	"github.com/nei7/gls/internal/repo"
	"github.com/nei7/gls/internal/service"
	"github.com/nei7/gls/pkg/video"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var env, addr string

	flag.StringVar(&env, "env", ".env", "Enviroment variables filename")
	flag.StringVar(&addr, "addr", ":3001", "Server address")

	flag.Parse()

	errC, err := run(env, addr)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running %s", err)
	}

}

func run(env, addr string) (<-chan error, error) {
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
	userRepo := repo.NewUserRepo(pool)

	videoService := service.NewVideoService(videoRepo)
	userService := service.NewUserService(logger, userRepo)

	tokenManager := service.NewTokenManager(viper.GetString("JWT_KEY"))

	videoServer := grpc_service.NewVideoServer(viper.GetString("VIDEO_STORAGE_PATH"), videoService, userService, tokenManager, logger)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		logger.Info("Shutdown signal received")

		defer func() {
			_ = logger.Sync()
			pool.Close()
			stop()
			close(errChan)
		}()

		if err := lis.Close(); err != nil {
			errChan <- err
		}
		logger.Info("Shutdown completed")

	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", addr))

		if err := newGRPCServer(*videoServer, lis); err != nil {
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
