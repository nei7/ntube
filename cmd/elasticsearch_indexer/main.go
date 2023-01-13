package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nei7/ntube/internal"
	"github.com/nei7/ntube/internal/datastruct"
	"github.com/nei7/ntube/internal/elasticsearch"
	"github.com/nei7/ntube/internal/kafka_service"
	"github.com/nei7/ntube/internal/opentelemetry"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	var env string

	flag.StringVar(&env, "env", ".env", "Enviroment variables filename")
	flag.Parse()

	errC, err := run(env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env string) (<-chan error, error) {
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

	var config kafka_service.KafkaConfig
	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	esClient, err := elasticsearch.NewElasticSearch()
	if err != nil {
		return nil, err
	}

	kafka, err := kafka_service.NewKafkaConsumer(config, "eslasticsearch_indexer")
	if err != nil {
		return nil, err
	}

	shutdown, err := opentelemetry.InitProviderWithJaegerExporter("eslasticsearch_indexer")
	if err != nil {
		return nil, err
	}

	srv := Server{
		logger: logger,
		kafka:  kafka,
		video:  elasticsearch.NewElasticVideo(esClient),
		doneC:  make(chan struct{}),
		closeC: make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			_ = logger.Sync()

			_ = kafka.Unsubscribe()

			stop()
			cancel()
			shutdown(ctxTimeout)
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}

	}()

	return errC, nil
}

type Server struct {
	logger *zap.Logger
	kafka  *kafka.Consumer
	video  *elasticsearch.Video
	doneC  chan struct{}
	closeC chan struct{}
}

func (s *Server) ListenAndServe() error {
	commit := func(msg *kafka.Message) {
		if _, err := s.kafka.CommitMessage(msg); err != nil {
			s.logger.Error("commit failed", zap.Error(err))
		}
	}

	go func() {
		for {
			select {
			case <-s.closeC:
				break

			default:
				msg, ok := s.kafka.Poll(150).(*kafka.Message)
				if !ok {
					continue
				}

				var evt struct {
					Type  string
					Value datastruct.Video
				}

				if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&evt); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					commit(msg)
					continue
				}

				ok = false

				switch evt.Type {
				case kafka_service.VideoEventCreated:
					if err := s.video.Index(context.Background(), evt.Value); err == nil {
						ok = true
					}
				case kafka_service.VideoEventDeleted:
					if err := s.video.Delete(context.Background(), evt.Value.ID); err == nil {
						ok = true
					}
				}

				if ok {
					s.logger.Info("Consumed", zap.String("type", evt.Type))
				}
			}
		}

		s.logger.Info("No more messages to consume. Exiting")
		s.doneC <- struct{}{}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shuting down server")
	close(s.closeC)

	for {

		select {
		case <-ctx.Done():
			return internal.WrapErrorf(ctx.Err(), internal.ErrorCodeUnknown, "contex.Done")
		case <-s.doneC:
			return nil
		}

	}
}
