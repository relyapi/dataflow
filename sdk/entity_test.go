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
	// [sink] item: {'data': {'name': 'xiaoming', 'age': 25, 'time': '2025-07-21 15:46:01', 'hello': 'world', 'world': 1121211}, 'source': '', 'url': 'https://www.json.cn9/', 'crawl_time': '2025-07-21 15:46:01'}
	record := Record{
		Data: map[string]any{
			"name": "opentome",
		},
		Metadata: Metadata{
			Url: "https://linux.do/t/topic/73111",
		},
	}

	if err := rs.SaveItem(record); err != nil {
		log.Fatal("发送失败:", err)
	}
}
