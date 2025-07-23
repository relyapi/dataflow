package model

import (
	"encoding/json"
)

type SinkType int

const (
	RAW     SinkType = iota // 日志、列表等原始数据
	ITEM                    // 清洗好的数据
	COMMENT                 // 评论数据
	PROFILE                 // 用户信息
)

// Resource 存储到数据库的字段
type Resource struct {
	// 基础元信息
	// sink信息（关联sink表）
	SinkId string `json:"sink_id"`
	// sink类型
	SinkType int `json:"sink_type"`

	// 来源
	ParentUrl string `json:"parent_url"`

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
