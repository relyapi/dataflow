package sink

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlSink struct {
	db        *gorm.DB
	tableName string
	log       *log.Helper
}

func (mysqlSink *MysqlSink) Sink(resources []*utils.Resource) (err error) {
	return mysqlSink.db.Table(mysqlSink.tableName).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "store_id"}},
		// 一旦创建不允许更新 sink_id和sink_type
		DoUpdates: clause.AssignmentColumns([]string{"store_key", "data", "metadata", "crawl_time"}),
	}).Create(&resources).Error
}
