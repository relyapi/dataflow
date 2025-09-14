package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/tomeai/dataflow/internal/sink"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSourceUseCase, NewSinkService)

func NewSinkService(itemSvc ItemRepo, logger log.Logger) *SinkService {
	loggerHelper := log.NewHelper(log.With(logger, "module", "pipeline.proto/sink"))
	sinkService := &SinkService{
		log: loggerHelper,
	}

	// 从数据库初始化
	sinkResults, err := itemSvc.GetAllSink(context.Background())
	if err != nil {
		loggerHelper.Fatalf("GetSink err: %v", err)
	}
	for _, item := range sinkResults {
		loggerHelper.Info("load sink: ", item)
		if item.SinkId != "" {
			switch item.Source.Type {
			case "mysql":
				dataSink, err := sink.NewMysqlSink(item, logger)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			case "mongo":
				dataSink, err := sink.NewMongoSink(item, logger)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			case "zincsearch":
				dataSink, err := sink.NewZincSearchSink(item, logger)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			case "cos":
				// 腾讯云对象存储
				dataSink, err := sink.NewCosSink(item, logger)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			default:
				loggerHelper.Info("not support sink type")
			}
		}
	}
	return sinkService
}
