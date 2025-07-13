package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/tomeai/dataflow/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewDataServiceManager, NewConfigService)

func NewDataServiceManager(sinkService *biz.SinkService, logger log.Logger) *DataServiceManager {
	return &DataServiceManager{
		sinkSvc: sinkService,
		log:     log.NewHelper(log.With(logger, "module", "service/data")),
	}
}

func NewConfigService(logger log.Logger) *ConfigServiceManager {
	return &ConfigServiceManager{
		log: log.NewHelper(log.With(logger, "module", "service/config")),
	}
}
