package service

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/tomeai/dataflow/api/v1/sink"
	"github.com/tomeai/dataflow/internal/biz"
	"github.com/tomeai/dataflow/internal/utils"
	"io"
	"strings"
)

// DataServiceManager 数据接收服务
type DataServiceManager struct {
	sink.UnimplementedDataHubServer
	sinkSvc *biz.SinkService
	log     *log.Helper
}

func (dataSvc *DataServiceManager) deserialize(msg *sink.DoSinkRequest) (data []*utils.Resource, err error) {
	if err := sonic.Unmarshal(msg.GetData(), &data); err != nil {
		return data, err
	}

	// 数据修正
	for _, item := range data {
		if item != nil {
			item.SinkId = msg.GetSinkId()
			if item.Url == "" {
				item.UrlMd5 = strings.ReplaceAll(uuid.New().String(), "-", "")
			} else {
				item.UrlMd5 = utils.CalcMD5(item.Url)
			}
		}
	}
	return data, nil
}

func (dataSvc *DataServiceManager) handleData(msg *sink.DoSinkRequest) error {
	if msg.GetSinkId() == "" {
		return errors.New("sink id is required")
	}

	// 必须是json
	data, err := dataSvc.deserialize(msg)
	if err != nil {
		return err
	}
	return dataSvc.sinkSvc.InsertData(msg.GetSinkId(), data)
}

func (dataSvc *DataServiceManager) DoSink(stream sink.DataHub_DoSinkServer) error {
	// 需要根据模板id进行解析并入库  校验也需要添加
	for {
		records, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&sink.Response{
				Status: true,
			})
		}
		if err != nil {
			if strings.HasSuffix(err.Error(), "context canceled") {
				return nil
			}
			continue
		}
		err = dataSvc.handleData(records)
		if err != nil {
			dataSvc.log.Error(err.Error())
			return err
		}
	}
}
