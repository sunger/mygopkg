package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBService is a database engine object.
type DBService struct {
	Default *gorm.DB            // the default database engine
	List    map[string]*gorm.DB // database engine list
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

var dbService = &DBService{
	List: map[string]*gorm.DB{},
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
			Db = dbService.Default
			fmt.Println("配置了默认数据库")
		}
	}()
	err := loadDBConfig(list)
	if err != nil {
		fmt.Println("[gorm]" + err.Error())
		return
	}

	// logs.Debug(dbConfigs)

	for _, conf := range dbConfigs {
		if !conf.Enable {
			continue
		}
		fmt.Println("22222222222222")
		fmt.Println(conf.Driver)
		fmt.Println(conf.Connstring)

		if (conf.Driver == "sqlite3" || conf.Driver == "sqlite") && !FileExists(conf.Connstring) {
			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
			f, err := os.Create(conf.Connstring)
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
			fmt.Println("333333333" + conf.Connstring)
			engine, err = gorm.Open(sqlite.Open(conf.Connstring), config)
		} else if dft == "mysql" {
			engine, err = gorm.Open(mysql.Open(conf.Connstring), config)
		} else if dft == "postgres" {
			engine, err = gorm.Open(postgres.Open(conf.Connstring), config)
		}

		if err != nil {
			fmt.Println("[gorm]" + err.Error())
			errs = append(errs, err.Error())
			continue
		}
		//		engine.SetLogger(faygo.NewLog())
		// engine.LogMode(true)

		db, _ := engine.DB()

		db.SetMaxOpenConns(conf.MaxOpenConns)
		db.SetMaxIdleConns(conf.MaxIdleConns)

		dbService.List[conf.Name] = engine
		if conf.IsDefault {
			dbService.Default = engine
		}
	}

}

func GetAllDbs() *DBService {
	return dbService
}

//var dbService = func() (serv *DBService) {
//	serv = &DBService{
//		List: map[string]*gorm.DB{},
//	}
//
//	var errs []string
//	defer func() {
//		if len(errs) > 0 {
//			panic("[gorm] " + strings.Join(errs, "\n"))
//		}
//		if serv.Default == nil {
//			serv.Default = Db
//			fmt.Println("[gorm] the `default` 数据库必须配置启用")
//		}
//	}()
//
//	dbconn:=&DbConn{}
//
//	err := loadDBConfig(dbconn.List())
//	if err != nil {
//		fmt.Println("[gorm]" + err.Error())
//		return
//	}
//
//	// logs.Debug(dbConfigs)
//
//	for _, conf := range dbConfigs {
//		if !conf.Enable {
//			continue
//		}
//		fmt.Println(conf.Driver)
//		fmt.Println(conf.Connstring)
//
//		var engine *gorm.DB
//
//		dft := conf.Driver
//		if dft == "sqlite" {
//			engine, err = gorm.Open(sqlite.Open(conf.Connstring), Db.Config)
//		} else if dft == "mysql" {
//			engine, err = gorm.Open(mysql.Open(conf.Connstring), Db.Config)
//		} else if dft == "postgres" {
//			engine, err = gorm.Open(postgres.Open(conf.Connstring), Db.Config)
//		}
//
//		if err != nil {
//			fmt.Println("[gorm]" + err.Error())
//			errs = append(errs, err.Error())
//			continue
//		}
//		//		engine.SetLogger(faygo.NewLog())
//		// engine.LogMode(true)
//
//		db,_ :=engine.DB()
//
//		db.SetMaxOpenConns(conf.MaxOpenConns)
//		db.SetMaxIdleConns(conf.MaxIdleConns)
//
//		if conf.Driver == "sqlite3" && !FileExists(conf.Connstring) {
//			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
//			f, err := os.Create(conf.Connstring)
//			if err != nil {
//				fmt.Println("[gorm]" + err.Error())
//				errs = append(errs, err.Error())
//			} else {
//				f.Close()
//			}
//		}
//
//		serv.List[conf.Name] = engine
//		if conf.IsDefault {
//			serv.Default = engine
//		}
//	}
//	return
//}()
