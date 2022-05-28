package db

import (
	"fmt"
	"github.com/sunger/mygopkg/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
	"testing"
)

func setUpAll() func() {
	//配置文件目录
	path := "F:/sh/mygopkg/config"
	config.Init("test", path)
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //logger.Silent
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // love表将是love，不再是loves，即可成功取消表明被加s
		},
	}
	//初始化默认数据库
	InitDb(cfg)


	return func() {
		// 如果测试结果符合预期的输出信息
		fmt.Println("数据库初始化正常.")
	}
}


func setMysqlUpAll() func() {
	//配置文件目录
	path := "F:/sh/mygopkg/config"
	config.Init("testmysql", path)
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //logger.Silent
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // love表将是love，不再是loves，即可成功取消表明被加s
		},
	}
	//初始化默认数据库
	InitDb(cfg)

	return func() {
		// 如果测试结果符合预期的输出信息
		fmt.Println("数据库初始化正常.")
	}
}




//func TestInitDatabase(t *testing.T) {
//	tearDownAll := setUpAll()
//	tearDownAll()   // you cannot use defer tearDownAll()
//}



//
func TestInitSqliteDbConns(t *testing.T) {

	tearDownAll := setUpAll()
	Db.AutoMigrate(
		&DbConn{},
	)
	db := DbConn{}
	//先清空数据
	Db.Exec("delete from "+db.TableName())
	mainDbPath := config.GetConfig().GetString("sqlite.name")

	dbPath:=config.ParentDir(mainDbPath)

	for i := 0; i < 10; i++ {
		db1 := DbConn{}

		db1.Name = "sqlite数据库"+strconv.Itoa(i)
		db1.DbName =  "db"+strconv.Itoa(i)+".db"
		db1.DbDir = dbPath
		db1.Enable = 1
		db1.Driver = "sqlite"
		db1.Insert(strconv.Itoa(i))
	}

	//将以上链接方式加载到内存
	InitDbConns()



	alldb := GetAllDbs()

	len := len(alldb.List)

	for i, i2 := range alldb.List {
		t.Log("根据记录id获取1111：",i,i2)
	}

	t.Log("内存中的数据库链接数量：",len)
	for i := 0; i < 10; i++ {
		//dbName :=  "db"+strconv.Itoa(i)+".db"
		//_,exist :=DB(dbName)
		//t.Log("根据数据库名称获取：",dbName,exist)

		db, exist2 := DB(strconv.Itoa(i))
		t.Log("根据记录id获取：",strconv.Itoa(i),exist2)

		//分别迁移这10个数据库
		db.AutoMigrate(
			&DbConn{},
		)
		db1 := DbConn{}

		db1.Name = "sqlite数据库"
		db1.DbName =  "db"+strconv.Itoa(i)+".db"
		db1.DbDir = dbPath
		db1.Enable = 1
		db1.Driver = "sqlite"
		db.Create(db1)
	}


	tearDownAll()   // you cannot use defer tearDownAll()
}


//
func TestInitMysqlDbConns(t *testing.T) {
	return
	tearDownAll := setMysqlUpAll()
	//Db.Raw("CREATE DATABASE IF NOT EXISTS test")
	Db.AutoMigrate(
		&DbConn{},
	)
	db := DbConn{}
	//先清空数据
	Db.Exec("delete from " + db.TableName())

	for i := 0; i < 10; i++ {
		db1 := DbConn{}

		db1.Name = "mysql数据库"+strconv.Itoa(i)
		db1.DbName =  "dbmysql"+strconv.Itoa(i)
		db1.User = "root"
		db1.Pwd = "root"
		db1.Port = "3306"
		db1.Host = "localhost"
		db1.Enable = 1
		db1.Driver = "mysql"
		db1.Insert(strconv.Itoa(i))
	}

	//将以上链接方式加载到内存
	InitDbConns()



	alldb := GetAllDbs()

	len := len(alldb.List)

	t.Log("内存中的数据库链接数量：",len)
	for i := 0; i < 10; i++ {
		//dbName :=  "db"+strconv.Itoa(i)+".db"
		//_,exist :=DB(dbName)
		//t.Log("根据数据库名称获取：",dbName,exist)

		db,exist2 := DB(strconv.Itoa(i))
		t.Log("根据记录id获取：",strconv.Itoa(i),exist2)
		dbName := "dbmysql"+strconv.Itoa(i)

		//分别迁移这10个数据库
		db.AutoMigrate(
			&DbConn{},
		)

		db1 := DbConn{}
		db1.Name = "mysql数据库"+strconv.Itoa(i)
		db1.DbName = dbName
		db1.User = "root"
		db1.Pwd = "root"
		db1.Port = "3306"
		db1.Host = "localhost"
		db1.Enable = 1
		db1.Driver = "sqlite"
		db.Create(db1)
	}


	tearDownAll()   // you cannot use defer tearDownAll()
}
