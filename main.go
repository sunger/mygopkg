package main

import (
	"fmt"
	"github.com/sunger/mygopkg/config"
	"github.com/sunger/mygopkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	fmt.Println("1111111111")
	config.Init("development", "F:\\go\\mygopkg\\config")

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),//logger.Silent
	}
	//初始化默认数据库
	db.InitDb(cfg)
	//加载所有数据库
	db.LoadAllDbs()
	//初始化数据库连接表
	db.Db.AutoMigrate(
		&db.DbConn{},
	)

	db.GetAllDbs()

}
