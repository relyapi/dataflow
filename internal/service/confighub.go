package service

import (
	"data-flow/api/v1/config"
	"github.com/go-kratos/kratos/v2/log"
)

type ConfigServiceManager struct {
	config.UnimplementedConfigHubServer
	log *log.Helper
}
