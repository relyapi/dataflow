//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/tomeai/dataflow/internal/biz"
	"github.com/tomeai/dataflow/internal/conf"
	"github.com/tomeai/dataflow/internal/data"
	"github.com/tomeai/dataflow/internal/server"
	"github.com/tomeai/dataflow/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
