package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/tomeai/dataflow/api/v1/flow"
	"github.com/tomeai/dataflow/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSourceRepo)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, error) {
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if conf.Database.AutoMigrate {
		err = db.AutoMigrate(&flow.Sink{}, &flow.Source{})
		if err != nil {
			return nil, err
		}
	}

	d := &Data{
		db:  db,
		log: log.NewHelper(log.With(logger, "module", "flow/data")),
	}
	return d, nil
}
