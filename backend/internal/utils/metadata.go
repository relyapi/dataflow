package utils

import (
	"reflect"
	"time"
)

func metadataToMap(meta Metadata) map[string]interface{} {
	if meta.CrawlTime == "" {
		meta.CrawlTime = time.Now().Format("2006-01-02 15:04:05")
	}
	result := make(map[string]interface{})
	v := reflect.ValueOf(meta)
	t := reflect.TypeOf(meta)

	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Tag.Get("json")
		if key == "" {
			key = t.Field(i).Name
		}
		result[key] = v.Field(i).Interface()
	}
	return result
}

type Metadata struct {
	Url       string `json:"url"`
	CrawlTime string `json:"crawl_time"`
}
