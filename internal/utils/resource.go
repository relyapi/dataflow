package utils

import (
	"encoding/json"
	"github.com/tomeai/dataflow/api/v1/sink"
)

type Resource struct {
	SinkId    string          `json:"sink_id"`
	SinkType  sink.SinkType   `json:"sink_type"`
	Source    string          `json:"source"`
	UrlMd5    string          `json:"url_md5"`
	Url       string          `json:"url"`
	CrawlTime string          `json:"crawl_time"`
	Data      json.RawMessage `json:"data"`
}
