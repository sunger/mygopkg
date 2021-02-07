package db

import (
	"github.com/sunger/mygopkg/log"
	"strings"
	"fmt"
)

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



//settings表转map
//{"cms":{"cms.filewhw":"local"}}
var sets = make(map[string]map[string]string,0)
//根据key获取参数
func GetSet(key string) string {
	ss := strings.Split(key, ".")

	if len(ss) < 2 {
		log.GetLog().Error("调用GetSet方法必须提供的key中包含。"+key)
		return  ""
	}
	// 如果已经加载过了，直接返回
	if _, ok := sets[ss[0]]; ok {
		return sets[ss[0]][key]
	}else{
		mp := GetSets(ss[0]) //数据库中可能不存在配置集合
		if mp == nil {
			return  ""
		}
		return  mp[key]
	}

}
//根据key获取参数集合
func GetSets(key string) map[string]string {
	ss := strings.Split(key, ".")

	if len(ss) > 1 {
		log.GetLog().Error("调用GetSets方法必须提供的key中不能包含。"+key)
		return  nil
	}

	if v, ok := sets[key]; ok {
		return  v
	}else {
		//从数据库中加载key

		mp := getDbMap(key)

		if mp != nil{
			sets[key] = mp
		}

		return mp
	}
}

func getDbMap(key string) map[string]string {
	s:= Settings{}
	list,_ := s.List(key)
	if len(list) == 0{
		log.GetLog().Error("数据库中不存在的集合"+key)
		return  nil
	}else {
		mp := make(map[string]string,0)
		//将集合加到全局变量中
		for _, value := range list {
			//fmt.Print(value, "\t")
			mp[value.Id] = value.Val
		}
		return  mp
	}
}

//根据key重新加载模块变量，如果是空，加载全部参数
func ReloadSet(key string)  {
	//全部
	if key=="" {
		newmp := make(map[string]map[string]string,0)

		for k,_ := range sets {
			mp:= getDbMap(k)
			if mp !=nil {
				newmp[k] = mp
			}

		}

		sets = newmp
	} else {
		//当前项，分集合和独立项
		ss := strings.Split(key, ".")

		//更新集合
		if len(ss) == 1 {
			mp:= getDbMap(key)
			if mp !=nil {
				sets[key] = mp
			}
		}else{ //更新集合中的一项
			s:= Settings{}
			item, err := s.Get(key)
			if err != nil {
				sets[ss[0]][key] = item.Val
			}
		}
	}
}

