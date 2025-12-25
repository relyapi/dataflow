package biz

import (
	"errors"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/internal/model"
	"github.com/tomeai/dataflow/internal/sink"
)

type SinkService struct {
	SinkMap sync.Map
	log     *log.Helper
}

func (sinkSvc *SinkService) InsertData(sinkId string, resources []*model.Resource) error {
	// 根据 sink_id 查询对应的库名和表名
	if value, ok := sinkSvc.SinkMap.Load(sinkId); ok {
		sinkInfo := value.(sink.DataHubSink)
		return sinkInfo.Sink(resources)
	}
	return errors.New("sink not found")
}
