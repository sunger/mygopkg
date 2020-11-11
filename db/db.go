package db

import (
	"log"
	"os"

	"github.com/sunger/mygopkg/cfg"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

// func initDB() {
// 	Db = initSqlite() // gormDB()
// }

func init() {
	Db = initSqlite() // gormDB()
}

func gormDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			Colorful: true,        // 彩色打印
		},
	)
	dsn := "root:root1234@tcp(127.0.0.1:3306)/casbin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	return db
}

func initMysql() *gorm.DB {
	dsn := "root:123123@tcp(localhost:3307)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	return db
}

func initSqlite() *gorm.DB {
	yml := cfg.NewYamlCfg()
	db, err := gorm.Open(sqlite.Open(yml.Sqlite.Name), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	// db.AutoMigrate(&models.Role{},&models.Routers{},&models.Tenant{},&models.Users{})

	return db
}
