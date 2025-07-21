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
	UrlMd5 string `json:"url_md5"`
	// 采集连接
	Url string `json:"url"`
	// 采集时间
	CrawlTime string `json:"crawl_time"`
	// 数据
	Data json.RawMessage `json:"data"`
}
