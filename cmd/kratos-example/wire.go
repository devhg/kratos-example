//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/devhg/kratos-example/internal/biz"
	"github.com/devhg/kratos-example/internal/conf"
	"github.com/devhg/kratos-example/internal/data"
	"github.com/devhg/kratos-example/internal/server"
	"github.com/devhg/kratos-example/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *zap.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
