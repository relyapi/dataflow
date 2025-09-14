package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/api/v1/config"
)

type ConfigServiceManager struct {
	config.UnimplementedConfigHubServer
	log *log.Helper
}
