package db

import (
	"errors"
	"fmt"
	"github.com/shenyisyn/goft-gin/goft"
	"github.com/sunger/mygopkg/log"
	"github.com/sunger/mygopkg/model"
	"os"
	"path/filepath"
	syslog "log"
	"github.com/sunger/mygopkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// //连接集合，用户根据key获取对应的db（此集合根据数据库中的表生成）(主要满足一个web根据不同的配置，使用不同的数据库)
// var dbs map[string]*gorm.DB

//默认连接，根据配置文件中获取的连接地址
var Db *gorm.DB

//根据配置文件初始化数据库
func InitDb(cfg *gorm.Config) {
	// config.Init(name)
	c := config.GetConfig()
	fmt.Println("配置文件初始化,database.default:" + c.GetString("database.default"))
	if c == nil {
		fmt.Println("配置文件未初始化，数据库未初始化")
		return
	}

	dft := c.GetString("database.default")
	if dft == "sqlite" {
		dir := c.GetString("sqlite.dir")
		name := c.GetString("sqlite.name")
		Db = InitSqlite(sqliteConn(dir, name), cfg)
	} else if dft == "mysql" {
		// "root:root1234@tcp(127.0.0.1:3306)/casbin?charset=utf8mb4&parseTime=True&loc=Local"
		user := c.GetString("mysql.user")
		password := c.GetString("mysql.password")
		host := c.GetString("mysql.host")
		port := c.GetString("mysql.port")
		name := c.GetString("mysql.name")
		dsn := mysqlConn(user, password, host, port, name)
		Db = InitMysql(dsn, cfg)
	} else if dft == "postgres" {
		// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		user := c.GetString("postgres.user")
		password := c.GetString("postgres.password")
		host := c.GetString("postgres.host")
		port := c.GetString("postgres.port")
		name := c.GetString("postgres.name")
		dsn := postgresConn(user, password, host, port, name)
		Db = InitPostgres(dsn, cfg)
	}

}

// //根据数据库中的记录初始化数据库集合
// func InitDbs() {
// 	loadDBConfig()
// }

// //根据key获取db连接对象
// func GetDb(key string) *gorm.DB {
// 	if value, ok := dbs[key]; ok {
// 		return value
// 	} else {
// 		return nil
// 	}
// }

func gormDB() *gorm.DB {
	newLogger := logger.New(
		syslog.New(os.Stdout, "\r\n", syslog.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			Colorful: true,        // 彩色打印
		},
	)
	dsn := "root:root1234@tcp(127.0.0.1:3306)/casbin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	return db
}

func InitPostgres(dsn string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	postgresDB, err := db.DB()
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	postgresDB.SetMaxIdleConns(5)
	postgresDB.SetMaxOpenConns(10)
	return db
}

func InitMysql(dsn string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	mysqlDB, err := db.DB()
	if err != nil {
		log.GetLog().Error(err.Error())
	}
	mysqlDB.SetMaxIdleConns(5)
	mysqlDB.SetMaxOpenConns(10)
	return db
}

func InitSqlite(name string, cfg *gorm.Config) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(name), cfg)

	if err != nil {
		log.GetLog().Error(err.Error())
	}
	// db.AutoMigrate(&models.Role{},&models.Routers{},&models.Tenant{},&models.Users{})

	return db
}

func sqliteConn(dir, name string) string {
	return filepath.ToSlash(filepath.Join(dir, name))
}

func mysqlConn(user, password, host, port, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
}
func postgresConn(user, password, host, port, name string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, name, port)
}

// 获取实体对应的db对象
func GetDb(m *model.BModel)  *gorm.DB {

	if m.DbKey == "" {
		log.GetLog().Error("数据库没有配置")
		goft.Error(errors.New("数据库没有配置"))
		return nil
	}
	db_,ok := DB(m.DbKey)
	if ok {
		log.GetLog().Info("数据库DbKey ="+m.DbKey + " dbname="+db_.Name())
		return db_
	}

	log.GetLog().Error("没有找到数据库连接")
	goft.Error(errors.New("没有找到数据库连接"))
	return nil
	//err = db.Db.Find(&results).Error
	//return results, err
}

