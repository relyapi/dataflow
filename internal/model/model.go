package model

import (
	"encoding/json"
)

type CrawlSource int

const (
	PC CrawlSource = iota
	ANDROID
	MINI // 小程序
	IOS
)

type CrawlType int

const (
	ITEM    CrawlType = iota // 清洗好的数据
	LIST                     // 列表等原始数据
	LOG                      // 日志数据
	COMMENT                  // 评论数据
	PROFILE                  // 用户信息
)

// Resource 存储到数据库的字段
type Resource struct {
	// 基础元信息
	// sink信息（关联sink表）
	SinkId string `json:"sink_id"`
	// 唯一标识
	// 数据库upsert唯一key,如果url为空，取uuid替代
	StoreId string `json:"store_id"`

	// 爬虫来源
	CrawlSource int `json:"crawl_source"`
	// 爬虫类型
	CrawlType int `json:"crawl_type"`
	// 来源 根据该链接生成md5 也是根据该url进行upsert resource只保留最新
	// 后续将通过数据同步 维护每次的版本变动
	CrawlUrl string `json:"crawl_url"`

	// 主域名 为了区分同一类下面的爬虫  比如 高校爬虫  huel.edu.cn tsing.edu.cn
	Hostname string `json:"hostname"`

	// 采集数据
	Data json.RawMessage `json:"data"`

	// 元数据
	Metadata json.RawMessage `json:"metadata"`

	// 采集时间
	CrawlTime string `json:"crawl_time"`
}
