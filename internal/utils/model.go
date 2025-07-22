package utils

import (
	"encoding/json"
)

// Resource 存储到数据库的字段
type Resource struct {
	// 基础元信息
	// sink信息（关联sink表）
	SinkId string `json:"sink_id"`

	// 唯一标识
	// 数据库upsert唯一key,如果url为空，取uuid替代
	StoreId string `json:"store_id"`
	// 采集连接 url或者key
	StoreKey string `json:"store_key"`

	// 采集数据
	Data json.RawMessage `json:"data"`

	// 元数据
	Metadata json.RawMessage `json:"metadata"`

	// 采集时间
	CrawlTime string `json:"crawl_time"`
}
