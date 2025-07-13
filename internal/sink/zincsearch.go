package sink

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/internal/utils"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"time"
)

type ZincSearchSink struct {
	zincDocument client.Document
	dbName       string
	tableName    string
	index        client.Index
	log          *log.Helper
}

func (zs *ZincSearchSink) convertData(data []*utils.Resource) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range data {
		// struct → []byte
		b, err := sonic.Marshal(item)
		if err != nil {
			continue // 或 log.Println(err)
		}

		// []byte → map
		var m map[string]interface{}
		if err := sonic.Unmarshal(b, &m); err != nil {
			continue
		}

		result = append(result, m)
	}
	return result
}

func (zs *ZincSearchSink) Sink(resources []*utils.Resource) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	query := client.NewMetaJSONIngest()
	query.SetIndex(zs.tableName)
	query.SetRecords(zs.convertData(resources))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, r, err := zs.zincDocument.Bulkv2(ctx).Query(*query).Execute()
	if err != nil {
		zs.log.Info(fmt.Sprintf("zinc response: %s", r.Status))
		return err
	}
	return nil
}
