package sink

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlSink struct {
	db        *gorm.DB
	tableName string
	log       *log.Helper
}

func (mysqlSink *MysqlSink) Sink(resources []*model.Resource) (err error) {
	return mysqlSink.db.Table(mysqlSink.tableName).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "store_id"}, {Name: "sink_id"}},
		// 一旦创建不允许更新 sink_id request_url hostname
		DoUpdates: clause.AssignmentColumns([]string{"sink_type", "data", "metadata", "crawl_time"}),
	}).Create(&resources).Error
}
