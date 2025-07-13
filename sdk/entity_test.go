package sdk

import (
	"github.com/tomeai/dataflow/api/v1/sink"
	"log"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	rs, err := NewResultService("77963b7a931377ad4ab5ad6a9cd718aa", sink.SinkType_INSERT)
	if err != nil {
		log.Fatal("初始化失败:", err)
	}

	record := Record{
		Item: map[string]any{
			"name": "opentome",
		},
		Metadata: Metadata{
			Source:    "list",
			Url:       "https://linux.do/t/topic/73111",
			CrawlTime: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	if err := rs.SaveItem(record); err != nil {
		log.Fatal("发送失败:", err)
	}
}
