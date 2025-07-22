package sdk

import (
	"log"
	"testing"
)

func TestSave(t *testing.T) {
	rs, err := NewResultService("77963b7a931377ad4ab5ad6a9cd718aa")
	if err != nil {
		log.Fatal("初始化失败:", err)
	}
	// [sink] item: {'data': {'name': 'xiaoming', 'age': 25, 'time': '2025-07-21 15:46:01', 'hello': 'world', 'world': 1121211}, 'url': 'https://www.json.cn9/', 'crawl_time': '2025-07-21 15:46:01'}
	record := Record{
		StoreKey: "https://linux.do/t/topic/73111111",
		Data: map[string]any{
			"name": "opentome",
		},
		Metadata: map[string]any{
			"name": "gage",
		},
	}

	if err := rs.SaveItem(record); err != nil {
		log.Fatal("发送失败:", err)
	}
}
