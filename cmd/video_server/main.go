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

	esv7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/elasticsearch"
	"github.com/nei7/ntube/internal/kafka_service"
	"github.com/nei7/ntube/internal/middlewares"
	"github.com/nei7/ntube/internal/opentelemetry"
	"github.com/nei7/ntube/internal/repo"
	"github.com/nei7/ntube/internal/rest"
	"github.com/nei7/ntube/internal/server"
	"github.com/nei7/ntube/internal/service"
	"github.com/nei7/ntube/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	esClient, err := elasticsearch.NewClient()
	if err != nil {
		return nil, err
	}

	config := kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("KAFKA_HOST"),
	}

	client, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	shutdown, err := opentelemetry.InitProviderWithJaegerExporter("video_server")
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		return nil, err
	}

	srv, err := newHttpServer(httpServerConfig{
		addr:          http_addr,
		logger:        logger,
		middlewares:   []func(next http.Handler) http.Handler{middlewares.LoggerMiddleware(*logger)},
		elasticSearch: esClient,
		kafka:         client,
		pool:          pool,
	})
	if err != nil {
		return nil, err
	}

	grpcSrv := newGRPCServer(grpcServerConfig{
		logger:        logger,
		kafka:         client,
		pool:          pool,
		elasticSearch: esClient,
	})

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
			defer shutdown(ctx)
			client.Close()
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
		logger.Info("Listening and serving", zap.String("grpc_address", grpc_addr))

		if err := grpcSrv.Serve(lis); err != nil {
			errChan <- err
		}

	}()

	go func() {
		logger.Info("Listening and serving", zap.String("http_address", http_addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return nil, nil
}

type grpcServerConfig struct {
	logger        *zap.Logger
	kafka         *kafka.Producer
	pool          *pgxpool.Pool
	elasticSearch *esv7.Client
}

func newGRPCServer(conf grpcServerConfig) *grpc.Server {
	videoRepo := repo.NewVideRepo(conf.pool)

	msgBroker := kafka_service.NewVideo(conf.kafka, viper.GetString("KAFKA_TOPIC"))
	search := elasticsearch.NewVideo(conf.elasticSearch)

	videoService := service.NewVideoService(conf.logger, videoRepo, search, msgBroker)

	tokenManager := service.NewTokenManager(viper.GetString("JWT_KEY"))
	ffmpegService := service.NewFfpmegService()

	videoUpload := service.NewVideoUpload(viper.GetString("VIDEO_STORAGE_PATH"), ffmpegService, videoService)

	videoServer := server.NewVideoServer(videoUpload, tokenManager, conf.logger)

	grpcServer := grpc.NewServer()
	pkg.RegisterVideoUploadServiceServer(grpcServer, videoServer)

	return grpcServer
}

type httpServerConfig struct {
	addr          string
	logger        *zap.Logger
	middlewares   []func(next http.Handler) http.Handler
	elasticSearch *esv7.Client
	pool          *pgxpool.Pool
	kafka         *kafka.Producer
}

func newHttpServer(conf httpServerConfig) (*http.Server, error) {
	router := chi.NewRouter()

	for _, mw := range conf.middlewares {
		router.Use(mw)
	}

	search := elasticsearch.NewVideo(conf.elasticSearch)
	videoRepo := repo.NewVideRepo(conf.pool)
	msgBroker := kafka_service.NewVideo(conf.kafka, viper.GetString("KAFKA_TOPIC"))

	svc := service.NewVideoService(conf.logger, videoRepo, search, msgBroker)
	rest.NewVideoHandler(viper.GetString("VIDEO_STORAGE_PATH"), svc).Register(router)
	router.Handle("/metrics", promhttp.Handler())

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
