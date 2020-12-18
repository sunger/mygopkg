package db

import (
	"fmt"
	"github.com/sunger/mygopkg/log"
	"github.com/sunger/mygopkg/model"
	"go.uber.org/zap"
)

// 数据库链接配置表
type DbConn struct {

	//名称
	Name string `gorm:"column:name;size:50"`
	//数据库名称
	DbName string `gorm:"column:dbname;size:50"`

	//Host
	Host string `gorm:"column:host;size:20"`

	//Port
	Port string `gorm:"column:port;size:10"`

	//数据库类型
	Driver string `gorm:"column:driver;size:20"`
	//账号
	User string `gorm:"column:user;size:20"`
	//如果是sqlite，数据库目录
	DbDir string `gorm:"column:dbdir;size:250"`
	//密码
	Pwd string `gorm:"column:pwd;size:30"`
	//数据库地址
	//Connstring string `gorm:"column:conn;size:250"`
	MaxOpenConns int `gorm:"maxopenconns"`
	MaxIdleConns int `gorm:"maxidleconns"`
	//默认,没有其他更小范围的连接配置，使用此连接
	IsDefault int `gorm:"column:isdft;size:1"`
	LogLevel   int `gorm:"column:loglv;size:1"`
	//是否可用
	Enable int `gorm:"column:enable;size:1"`

	model.BModel
}

func (DbConn) TableName() string {
	return "s_dbconn"
}

func (u *DbConn) Insert(newid string) (id string) {
	if newid=="" || len(newid) == 0{
		u.CreateId()
	}else{
		u.Id = newid
	}

	Db.Create(&u)
	return u.Id
}

func (apps *DbConn) Del(id string) (err error) {
	return Db.Where("id = ?", id).Delete(DbConn{}).Error
}

// 更新实体
func (apps *DbConn) Update() (err error) {
	a := DbConn{}
	a.Id = apps.Id
	return Db.Model(&a).Updates(map[string]interface{}{
		"name":   apps.Name,
		"dbname": apps.DbName,
		"host": apps.Host,
		"driver": apps.Driver,
		"port": apps.Port,
		"dbdir":apps.DbDir,
		"enable":apps.Enable,
		"maxopenconns":apps.MaxOpenConns,
		"maxidleconns":apps.MaxIdleConns,
		"loglv":apps.LogLevel,
		"user":   apps.User,
		"pwd":    apps.Pwd,
		"isdft":  apps.IsDefault}).Error

}

//根据主键获取实体
func (apps *DbConn) Get(id string) (DbConn, error) {
	a := DbConn{}
	a.Id = id
	err := Db.Where(&a).Find(&a).Error
	return a, err
}

//分页方法
func (b *DbConn) List() (list []DbConn) {
	Db.Find(&list)
	return  list
}


//分页方法
func (b *DbConn) PageList(page, size int, filter string, sort string) ([]DbConn, int) {

	table := b.TableName()

	if size == 0 {
		size = 20
	}

	var offset int
	if page <= 0 {
		offset = 0
	} else {
		offset = page * size
	}

	if len(filter) == 0 {
		filter = "1=1"
	}
	if len(sort) == 0 {
		sort = "id desc"
	}

	arrSql := fmt.Sprintf("SELECT * FROM %s WHERE %s order by %s limit %d,%d",
		table, filter, sort, offset, size)

	countSql := fmt.Sprintf("SELECT count(0) as total FROM %s WHERE %s",
		table, filter)

	results := make([]DbConn, 0)

	rows, err := Db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		log.GetLog().Error("会员分页错误", zap.String("", err.Error()))
	}

	for rows.Next() {
		var a DbConn
		Db.ScanRows(rows, &a)
		results = append(results, a)
	}

	var total model.PageTotal
	Db.Raw(countSql).Scan(&total)
	count := total.Total

	if err != nil {
		fmt.Println(err)
	}

	return results, count
}
