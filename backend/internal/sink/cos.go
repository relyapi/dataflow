package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tomeai/dataflow/internal/model"
	"net/http"
	"sync"
	"time"
)

type CosSink struct {
	cosClient  *cos.Client
	tableName  string
	threadPool int
	log        *log.Helper
}

// Sink 启动并发上传
func (cosSink *CosSink) Sink(resources []*model.Resource) error {
	return cosSink.batchUpload(resources)
}

// 并发上传
func (cosSink *CosSink) batchUpload(resources []*model.Resource) error {
	resourceCh := make(chan *model.Resource, len(resources))
	var wg sync.WaitGroup

	// 启动 worker goroutines
	for i := 0; i < cosSink.threadPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range resourceCh {
				if err := cosSink.uploadResource(res); err != nil {
					cosSink.log.Errorf("Upload error: %v", err)
				}
			}
		}()
	}

	// 发送资源到 channel
	for _, res := range resources {
		resourceCh <- res
	}
	close(resourceCh)

	wg.Wait()
	return nil
}

// 实际上传逻辑
func (cosSink *CosSink) uploadResource(res *model.Resource) error {
	// 	key := "folder/" + res.FileName

	var metadataMap map[string]string

	if err := json.Unmarshal(res.Metadata, &metadataMap); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	header := &http.Header{}

	for key, value := range metadataMap {
		header.Add(fmt.Sprintf("x-cos-meta-%s", key), value)
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "application/octet-stream",
			XCosMetaXXX: header,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storeKey := fmt.Sprintf("%s/%s", res.Hostname, res.StoreId)
	_, err := cosSink.cosClient.Object.Put(ctx, storeKey, bytes.NewReader(res.Data), opt)
	return err
}
