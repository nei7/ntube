//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/nei7/ntube/app/2fa/internal/biz"
	"github.com/nei7/ntube/app/2fa/internal/conf"
	"github.com/nei7/ntube/app/2fa/internal/data"
	"github.com/nei7/ntube/app/2fa/internal/server"
	"github.com/nei7/ntube/app/2fa/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data_Database, *conf.Email, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
