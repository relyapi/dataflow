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
		Columns:   []clause.Column{{Name: "url_md5"}},
		DoUpdates: clause.AssignmentColumns([]string{"url", "source", "data", "crawl_time"}),
	}).Create(&resources).Error
}
