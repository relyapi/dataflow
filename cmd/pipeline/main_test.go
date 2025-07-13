package main

import (
	"context"
	"data-flow/api/v1/sink"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"log"
	"strings"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithTimeout(5*time.Second),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	client := sink.NewDataHubClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. 获取流
	stream, err := client.DoSink(ctx)
	if err != nil {
		log.Fatalf("打开 DoSink 流失败: %v", err)
	}

	// 2. 发送多条消息
	data := `[{"source":"list","url":"https://linux.do/t/topic/761213","data":{"name":"opentome1"},"crawl_time":"2025-07-13 13:56:12"},{"source":"list","url":"https://linux.do/t/topic/761213","data":[{"name":"opentome122"},{"name":"opentome222"}],"crawl_time":"2025-07-13 13:56:12"}]`
	req := &sink.DoSinkRequest{
		// 填充请求字段
		SinkId: "77963b7a931377ad4ab5ad6a9cd718aa",
		Data:   []byte(data),
	}

	if err := stream.Send(req); err != nil {
		log.Fatalf("发送数据失败: %v", err)
	}

	// 3. 关闭发送并接收服务端响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("关闭发送流失败: %v", err)
	}

	log.Printf("接收到响应: %+v", resp)
}

func TestEs(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://search-c1j.public.cn-beijing.es-serverless.aliyuncs.com:9200",
		},
		// 如果需要认证，可以添加：
		Username: "search-c1j",
		Password: "whEnLNxHVZoQ3xdGp4fMiIRClpzIyM",
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return
	}
	indexName := "test_index"
	// 定义索引映射
	mapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "integer"},
				"title": {"type": "text", "analyzer": "standard"},
				"content": {"type": "text", "analyzer": "standard"},
				"author": {"type": "keyword"},
				"tags": {"type": "keyword"}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		return
	}

	fmt.Printf("Index '%s' created successfully\n", indexName)

}
