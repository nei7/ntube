// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/nei7/ntube/app/user/internal/biz"
	"github.com/nei7/ntube/app/user/internal/conf"
	"github.com/nei7/ntube/app/user/internal/data"
	"github.com/nei7/ntube/app/user/internal/server"
	"github.com/nei7/ntube/app/user/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {

	pool, err := data.NewPgxPool(confData.Database)
	if err != nil {
		return nil, nil, err
	}

	dataData, cleanup, err := data.NewData(pool, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewUserRepo(dataData, logger)
	greeterUsecase := biz.NewUserUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
