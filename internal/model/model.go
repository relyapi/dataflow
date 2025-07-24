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

type CrawlSource int

const (
	PC      = 0
	ANDROID = 1 // 手机端
	IOS     = 2
	MINI    = 3 // 小程序
)

// Resource 存储到数据库的字段
type Resource struct {
	// 基础元信息
	// sink信息（关联sink表）
	SinkId string `json:"sink_id"`
	// sink类型
	SinkType int `json:"sink_type"`

	// 主域名 为了区分同一类下面的爬虫  比如 高校爬虫  huel.edu.cn tsing.edu.cn
	Hostname string `json:"hostname"`

	// 唯一标识
	// 数据库upsert唯一key,如果url为空，取uuid替代
	StoreId string `json:"store_id"`
	// 来源 根据该链接生成md5 也是根据该url进行upsert resource只保留最新
	// 后续将通过数据同步 维护每次的版本变动
	RequestUrl string `json:"request_url"`

	// 采集数据
	Data json.RawMessage `json:"data"`

	// 元数据
	Metadata json.RawMessage `json:"metadata"`

	// 采集时间
	CrawlTime string `json:"crawl_time"`
}
