package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/tomeai/dataflow/api/v1/config"
	"github.com/tomeai/dataflow/api/v1/sink"
	"github.com/tomeai/dataflow/internal/conf"
	service2 "github.com/tomeai/dataflow/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, data *service2.DataServiceManager, configServer *service2.ConfigServiceManager, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	sink.RegisterDataHubServer(srv, data)
	config.RegisterConfigHubServer(srv, configServer)

	return srv
}
