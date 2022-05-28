package db

import (
	"errors"
	"fmt"
	syslog "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunger/mygopkg/config"
	"github.com/sunger/mygopkg/goft"
	"github.com/sunger/mygopkg/log"
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
// DBService is a database engine object.
type DBService struct {
	Default *gorm.DB            // the default database engine
	List    map[string]*gorm.DB // database engine list
}


var dbService = &DBService{
	List: map[string]*gorm.DB{},
}


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
		mainDsn := mysqlConn(user, password, host, port, "mysql")
		Db = InitMysql(name,dsn,mainDsn, cfg)
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

	GetAllDbs().Default = Db



}

//加载数据库中所有的数据库连接到内存中
func InitDbConns()  {
	LoadAllDbs()
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

//dbname 要链接的数据库名称
//mainDsn 主链接(mysql链接)，如果没有数据库，自动创建数据库
func InitMysql(dbname, dsn, mainDsn string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		log.GetLog().Error(err.Error())
		log.GetLog().Error("不存在数据库: " + dbname +" 链接到mysql主数据库创建其他数据库 " + mainDsn)
		db1, err1 := gorm.Open(mysql.Open(mainDsn), cfg)
		if err1 != nil {
			log.GetLog().Error("链接主数据库错误: " + mainDsn)
		}

		db1.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
		log.GetLog().Error("执行了语句 CREATE DATABASE IF NOT EXISTS " + dbname)
		//先创建数据库之后再链接数据库
		db, err = gorm.Open(mysql.Open(dsn), cfg)
		if err != nil {
			log.GetLog().Error("链接到mysql中创建数据库之后，打开数据库失败: " + dsn)
		}
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
		db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
		//创建成功之后再打开数据库链接
		db, err = gorm.Open(sqlite.Open(name), cfg)
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
func GetDb(key string) *gorm.DB {

	if key == "" {
		log.GetLog().Error("数据库没有配置")
		goft.Error(errors.New("数据库没有配置"))
		return nil
	}
	db_, ok := DB(key)
	if ok {
		log.GetLog().Info("数据库DbKey =" + key + " dbname=" + db_.Name())
		return db_
	}

	log.GetLog().Error("没有找到数据库连接key:" + key)
	goft.Error(errors.New("没有找到数据库连接key:" + key))
	return nil
	//err = db.Db.Find(&results).Error
	//return results, err
}


func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

//加载所有数据库,这里在主模块中调用
func LoadAllDbs() {
	dbconn := &DbConn{}
	MapListToDBService(dbconn.List(), Db.Config)
}

//将数据库记录对象DbConn集合转map，这里在主模块之外的模块中调用
//将数据库中默认数据库赋值给 dbService.Default
func MapListToDBService(list []DbConn, config *gorm.Config) {
	var errs []string
	defer func() {
		if len(errs) > 0 {
			panic("[MapListToDBService] " + strings.Join(errs, "\n"))
		}
		if dbService.Default == nil {
			//dbService.Default = Db
			fmt.Println("未配置默认数据库")
		}else{
			//Db = dbService.Default
			fmt.Println("配置了默认数据库")
		}
	}()
	//err := loadDBConfig(list)
	//if err != nil {
	//	fmt.Println("[gorm]" + err.Error())
	//	return
	//}

	// logs.Debug(dbConfigs)

	for _, conf := range list {
		if conf.Enable == 0 {
			continue
		}

		if (conf.Driver == "sqlite3" || conf.Driver == "sqlite") && !FileExists(conf.DbDir) {
			os.MkdirAll(filepath.Dir(conf.DbDir), 0777)
			f, err := os.Create(conf.DbDir)
			if err != nil {
				fmt.Println("[gorm]" + err.Error())
				errs = append(errs, err.Error())
			} else {
				f.Close()
			}
		}

		var engine *gorm.DB

		dft := conf.Driver
		if dft == "sqlite" {

			//engine, err = gorm.Open(sqlite.Open(conf.Connstring), config)
			engine = InitSqlite(sqliteConn(conf.DbDir, conf.DbName), config)
			fmt.Println("333333333",engine)
		} else if dft == "mysql" {
			//fmt.Println("333333333mysql" + conf.Connstring)
			//engine, err = gorm.Open(mysql.Open(conf.Connstring), config)

			dsn := mysqlConn(conf.User, conf.Pwd, conf.Host, conf.Port, conf.DbName)
			mainDsn := mysqlConn(conf.User, conf.Pwd, conf.Host, conf.Port, "mysql")
			engine = InitMysql(conf.DbName,dsn,mainDsn, config)
		} else if dft == "postgres" {
			//fmt.Println("333333333postgres" + conf.Connstring)
			//engine, err = gorm.Open(postgres.Open(conf.Connstring), config)

			user := conf.User
			password := conf.Pwd
			host := conf.Host
			port := conf.Port
			name := conf.DbName
			dsn := postgresConn(user, password, host, port, name)
			engine = InitPostgres(dsn, config)
		}

		//if err != nil {
		//	fmt.Println("[gorm]" + err.Error())
		//	errs = append(errs, err.Error())
		//	continue
		//}
		// engine.SetLogger(faygo.NewLog())
		// engine.LogMode(true)

		//db, _ := engine.DB()
		//
		//db.SetMaxOpenConns(conf.MaxOpenConns)
		//db.SetMaxIdleConns(conf.MaxIdleConns)

		dbService.List[conf.Id] = engine
		if conf.IsDefault == 1 {
			dbService.Default = engine
		}
	}

}

func GetAllDbs() *DBService {
	return dbService
}

