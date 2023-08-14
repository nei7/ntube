// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nei7/ntube/app/auth/internal/biz"
	"github.com/nei7/ntube/app/auth/internal/conf"
	"github.com/nei7/ntube/app/auth/internal/data"
	"github.com/nei7/ntube/app/auth/internal/server"
	"github.com/nei7/ntube/app/auth/internal/service"
	"go.opentelemetry.io/otel/sdk/trace"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, data_Database *conf.Data_Database, email *conf.Email, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	conn, err := data.NewPgxPool(data_Database)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(conn, logger)
	if err != nil {
		return nil, nil, err
	}
	authRepo := data.NewEmailVerifyRepo(dataData, logger)
	authUsecase := biz.NewAuthUsecase(authRepo, logger)
	emailSenderUsecase := biz.NewEmailSenderUsecase(email)
	emailJobService := service.NewEmailJobService(authUsecase, emailSenderUsecase, confServer)
	kafkaServer := server.NewKafkaServer(confServer, logger, emailJobService, tracerProvider)
	authService := service.NewAuthService(authUsecase)
	httpServer := server.NewHTTPServer(confServer, authService, logger, tracerProvider)
	app := newApp(logger, kafkaServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
