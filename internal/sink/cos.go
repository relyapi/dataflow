package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tomeai/dataflow/internal/utils"
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
func (cosSink *CosSink) Sink(resources []*utils.Resource) error {
	return cosSink.batchUpload(resources)
}

// 并发上传
func (cosSink *CosSink) batchUpload(resources []*utils.Resource) error {
	resourceCh := make(chan *utils.Resource, len(resources))
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
func (cosSink *CosSink) uploadResource(res *utils.Resource) error {
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

	_, err := cosSink.cosClient.Object.Put(ctx, res.StoreKey, bytes.NewReader(res.Data), opt)
	return err
}
