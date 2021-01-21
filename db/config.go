package db

import "fmt"

// import (
// 	"os"
// 	"path/filepath"
// )

// DBConfig is database connection config
type DBConfig struct {
	Name         string `ini:"-"`
	Enable       bool   `ini:"enable" comment:"Enable the config section"`
	Driver       string `ini:"driver" comment:"mssql | odbc(mssql) | mysql | mymysql | postgres | sqlite3 | oci8 | goracle"`
	Connstring   string `ini:"connstring" comment:"Connect String"`
	MaxOpenConns int    `ini:"max_open_conns"`
	MaxIdleConns int    `ini:"max_idle_conns"`
	IsDefault    bool   `ini:"IsDefault" comment:"print sql"`
}

var (
	dbConfigs     = map[string]*DBConfig{}
	DefaultConfig = &DBConfig{
		Name:         "default",
		Driver:       "sqlite3",
		Connstring:   "test.db",
		MaxOpenConns: 100,
		MaxIdleConns: 100,
		IsDefault:    true,
	}
)

func loadDBConfig(list []DbConn) error {

	//从数据库里加载数据库配置
	//var hadDefaultConfig bool
	for _, v := range list {

		var dbConfig *DBConfig
		//if v.IsDefault == 1 {
		//	//dbConfig = DefaultConfig
		//	DefaultConfig = &DBConfig{
		//		Name:         v.Id,
		//		Driver:       v.Driver,
		//		Connstring:   conn,
		//		Enable:       v.Enable == 1,
		//		MaxOpenConns: v.MaxOpenConns,
		//		MaxIdleConns: v.MaxIdleConns,
		//		IsDefault:    v.IsDefault == 1,
		//	}
		//
		//	hadDefaultConfig = true
		//} else {

			dft := v.Driver
			conn := ""
			if dft == "sqlite" {
				conn = sqliteConn(v.DbDir, v.DbName)
			} else if dft == "mysql" {
				conn = mysqlConn(v.User, v.Pwd, v.Host, v.Port, v.DbName)
			} else if dft == "postgres" {
				conn = postgresConn(v.User, v.Pwd, v.Host, v.Port, v.DbName)
			}

			dbConfig = &DBConfig{
				Name:         v.Id,
				Driver:       v.Driver,
				Connstring:   conn,
				Enable:       v.Enable == 1,
				MaxOpenConns: v.MaxOpenConns,
				MaxIdleConns: v.MaxIdleConns,
				IsDefault:    v.IsDefault == 1,
			}
		//}

		dbConfigs[dbConfig.Name] = dbConfig
		fmt.Println(dbConfig.Name + dbConfig.Driver)
	}

	//if !hadDefaultConfig {
	//	*DefaultConfig = DBConfig{}
	//}

	return nil
}
