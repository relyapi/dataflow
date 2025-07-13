package biz

import (
	"context"
	"data-flow/internal/sink"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSourceUseCase, NewSinkService)

func NewSinkService(itemSvc ItemRepo, logger log.Logger) *SinkService {
	loggerHelper := log.NewHelper(log.With(logger, "module", "pipeline.proto/sink"))
	sinkService := &SinkService{
		log: loggerHelper,
	}

	// 从数据库初始化
	// todo: 抽离单独的配置服务  创建还原sink
	sinkResults, err := itemSvc.GetAllSink(context.Background())
	if err != nil {
		loggerHelper.Fatalf("GetSink err: %v", err)
	}
	for _, item := range sinkResults {
		loggerHelper.Info("load sink: ", item)
		if item.SinkId != "" {
			switch item.Source.Type {
			case "mysql":
				dataSink, err := sink.NewMysqlSink(item)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			case "mongo":
				dataSink, err := sink.NewMongoSink(item)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			case "zincsearch":
				dataSink, err := sink.NewZincSearchSink(item)
				if err != nil {
					loggerHelper.Error(err)
				}
				sinkService.SinkMap.Store(item.SinkId, dataSink)
			//case "elasticsearch":
			default:
				loggerHelper.Info("not support sink type")
			}
		}
	}
	return sinkService
}
