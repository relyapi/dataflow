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
	// [sink] item: {'data': {'name': 'xiaoming', 'age': 25, 'time': '2025-07-21 15:46:01', 'hello': 'world', 'world': 1121211}, 'source': '', 'url': 'https://www.json.cn9/', 'crawl_time': '2025-07-21 15:46:01'}
	record := Record{
		Data: map[string]any{
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
