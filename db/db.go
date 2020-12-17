package db

import (
	"fmt"
	"log"
	"os"

	"github.com/sunger/mygopkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDb(name string, cfg *gorm.Config) {
	// config.Init(name)
	c := config.GetConfig()

	if c == nil {
		fmt.Println("配置文件未初始化，数据库未初始化")
		return nil
	}

	dft := c.GetString("database.default")
	if dft == "sqlite" {
		name := c.GetString("sqlite.name")
		Db = initSqlite(name, cfg)
	} else if dft == "mysql" {
		// "root:root1234@tcp(127.0.0.1:3306)/casbin?charset=utf8mb4&parseTime=True&loc=Local"
		user := c.GetString("mysql.user")
		password := c.GetString("mysql.password")
		host := c.GetString("mysql.host")
		port := c.GetString("mysql.port")
		name := c.GetString("mysql.name")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
		Db = initMysql(dsn, cfg)
	}

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

func initMysql(dsn string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), cfg)
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

func initSqlite(name string, cfg *gorm.Config) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(name), cfg)

	if err != nil {
		log.Fatal(err)
	}
	// db.AutoMigrate(&models.Role{},&models.Routers{},&models.Tenant{},&models.Users{})

	return db
}
