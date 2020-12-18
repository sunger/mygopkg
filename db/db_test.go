package db

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestBsModel(t *testing.T) {

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	InitDb(cfg)

	// 如果测试结果符合预期的输出信息
	t.Log("数据库初始化正常.")
}
