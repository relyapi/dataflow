package sdk

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/tomeai/dataflow/api/v1/sink"
)

// Record 表示一条数据 + 元信息
type Record struct {
	Item     any               `json:"data"`     // 实际数据内容
	Metadata map[string]string `json:"metadata"` // 源、url、时间等
}

type ResultService struct {
	sinkID    string
	sinkType  sink.SinkType
	sinkStub  sink.DataHubClient
	batchSize int
}

var (
	once    sync.Once
	svc     *ResultService
	initErr error
)

// 初始化 gRPC 客户端单例
func getClient() (sink.DataHubClient, error) {
	var stub sink.DataHubClient
	once.Do(func() {
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("127.0.0.1:9000"),
			grpc.WithTimeout(5*time.Second),
			grpc.WithMiddleware(recovery.Recovery()),
		)
		if err != nil {
			initErr = err
			return
		}
		stub = sink.NewDataHubClient(conn)
	})
	return stub, initErr
}

func NewResultService(sinkID string, sinkType sink.SinkType) (*ResultService, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	return &ResultService{
		sinkID:    sinkID,
		sinkType:  sinkType,
		sinkStub:  client,
		batchSize: 100,
	}, nil
}

// SaveItem 写入单条记录
func (r *ResultService) SaveItem(item Record) error {
	return r.SaveItems([]Record{item})
}

// SaveItems 批量写入记录，自动按 batchSize 分批
func (r *ResultService) SaveItems(items []Record) error {
	if r.sinkID == "" {
		log.Println("[warn] sinkID is empty, skipping send.")
		return nil
	}

	var batch []map[string]any
	for _, item := range items {
		record := make(map[string]any)
		record["data"] = item.Item
		for k, v := range item.Metadata {
			record[k] = v
		}
		log.Printf("[sink] item: %+v\n", record)
		batch = append(batch, record)

		if len(batch) >= r.batchSize {
			if err := r.send(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	// 发送剩余
	if len(batch) > 0 {
		if err := r.send(batch); err != nil {
			return err
		}
	}
	return nil
}

// 实际发送方法（JSON 序列化 + DoSink）
func (r *ResultService) send(records []map[string]any) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := r.sinkStub.DoSink(ctx)
	if err != nil {
		return err
	}

	if err := stream.Send(&sink.DoSinkRequest{
		SinkId:   r.sinkID,
		SinkType: r.sinkType,
		Data:     data,
	}); err != nil {
		return err
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("[sink] send done: %+v\n", resp)
	return nil
}
