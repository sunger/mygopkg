package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sunger/mygopkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDb(cfg *gorm.Config) {
	// config.Init(name)
	c := config.GetConfig()
	fmt.Println("配置文件初始化,database.default:"+c.GetString("database.default"))
	if c == nil {
		fmt.Println("配置文件未初始化，数据库未初始化")
		return
	}

	dft := c.GetString("database.default")
	if dft == "sqlite" {
		name := c.GetString("sqlite.name")
		Db = InitSqlite(name, cfg)
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

func InitPostgres(dsn string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		log.Fatal(err)
	}
	postgresDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	postgresDB.SetMaxIdleConns(5)
	postgresDB.SetMaxOpenConns(10)
	return db
}

func InitMysql(dsn string, cfg *gorm.Config) *gorm.DB {
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

func InitSqlite(name string, cfg *gorm.Config) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(name), cfg)

	if err != nil {
		log.Fatal(err)
	}
	// db.AutoMigrate(&models.Role{},&models.Routers{},&models.Tenant{},&models.Users{})

	return db
}

func sqliteConn(dir, name string) string {
	return filepath.Join(dir, name)
}

func mysqlConn(user, password, host, port, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
}
func postgresConn(user, password, host, port, name string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, name, port)
}
