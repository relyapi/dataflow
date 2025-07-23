package sdk

import (
	"context"
	"encoding/json"
	"github.com/tomeai/dataflow/internal/model"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/tomeai/dataflow/api/v1/sink"
)

// Record 表示一条数据 + 元信息
type Record struct {
	SinkType  model.SinkType `json:"sink_type"`
	ParentUrl string         `json:"parent_url"`
	StoreKey  string         `json:"store_key"`
	Data      any            `json:"data"`
	Metadata  any            `json:"metadata"`
	CrawlTime string         `json:"crawl_time"`
}

type ResultService struct {
	sinkID    string
	sinkStub  sink.DataHubClient
	batchSize int
}

var (
	once    sync.Once
	initErr error
)

func getGRPCAddr() string {
	addr := os.Getenv("DATAFLOW_GRPC_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:9000"
	}
	return addr
}

// 初始化 gRPC 客户端单例
func getClient() (sink.DataHubClient, error) {
	var stub sink.DataHubClient
	once.Do(func() {
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint(getGRPCAddr()),
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

func NewResultService(sinkID string) (*ResultService, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	return &ResultService{
		sinkID:    sinkID,
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

	var batch []Record
	for _, item := range items {
		if item.CrawlTime == "" {
			item.CrawlTime = time.Now().Format("2006-01-02 15:04:05")
		}

		log.Printf("[sink] item: %+v\n", item)
		batch = append(batch, item)

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
func (r *ResultService) send(records []Record) error {
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
		SinkId: r.sinkID,
		Data:   data,
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
