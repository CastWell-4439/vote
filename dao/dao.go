package dao

import (
	"toupiao/config"
	"toupiao/logger"
)

func init() {
	if config.Err != nil {
		logger.Error(map[string]interface{}{"fail to connect": config.Err.Error()})

	}
	if config.DB.Error != nil {
		logger.Error(map[string]interface{}{"mysql error": config.DB.Error})
	}
}
