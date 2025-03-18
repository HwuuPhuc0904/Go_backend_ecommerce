package global

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/pkg/logger"
	setting "GOLANG/github.com/HwuuPhuc0904/backend-api/settings"
	"gorm.io/gorm"
)


var (
	Config setting.Config
	Logger *logger.LoggerZap
	DB *gorm.DB
)