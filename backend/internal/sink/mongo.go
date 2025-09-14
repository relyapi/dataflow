package sink

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoSink struct {
	database  *mongo.Database
	tableName string
	log       *log.Helper
}

func (ms *MongoSink) convertData(data []*model.Resource) []interface{} {
	var result []interface{}
	for _, item := range data {
		result = append(result, item)
	}
	return result
}

func (ms *MongoSink) Sink(resources []*model.Resource) (err error) {
	// 根据 sink_type 判断是新增还是更新
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	collection := ms.database.Collection(ms.tableName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = collection.InsertMany(ctx, ms.convertData(resources))
	if err != nil {
		return err
	}
	return nil
}
