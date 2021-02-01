package main

import (
	"fmt"
	"github.com/sunger/mygopkg/nsq"

	"github.com/sunger/mygopkg/config"
	"github.com/sunger/mygopkg/db"
	"github.com/sunger/mygopkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	fmt.Println("1111111111")
	config.Init("development", "F:\\go\\mygopkg\\config")

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //logger.Silent
	}
	//初始化默认数据库
	db.InitDb(cfg)

	//初始化数据库连接表
	db.Db.AutoMigrate(
		&db.DbConn{},
	)

	//模拟初始数据库连接
	// mockDbConn()

	//加载所有数据库
	db.LoadAllDbs()
	db.GetAllDbs()

	db1 := db.MustDB("1") //数据库中记录id
	db2 := db.MustDB("2") //数据库中记录id

	//模拟初始化数据库连接表
	db1.AutoMigrate(
		&db.DbConn{},
	)
	db2.AutoMigrate(
		&db.DbConn{},
	)

	fmt.Println(db1.Name())
	fmt.Println(db2.Name())

	p := model.PageParams{}

	arr := model.GetFlts(p)
	filterstr := arr[0]
	orderstr := arr[1]

	fmt.Println(filterstr)
	fmt.Println(orderstr)

}

func mockDbConn() {
	db1 := db.DbConn{}

	db1.Name = "sqlite数据库1"
	db1.DbName = "db1.db"
	db1.DbDir = "testdb"
	db1.Enable = 1
	db1.Driver = "sqlite"
	db1.Insert("1")

	db2 := db.DbConn{}
	db2.Name = "sqlite数据库2"
	db2.DbName = "db2.db"
	db2.Enable = 1
	db2.DbDir = "testdb"
	db2.Driver = "sqlite"
	db2.Insert("2")

	nsq.Pg = nsq.RqProgram{}
}
