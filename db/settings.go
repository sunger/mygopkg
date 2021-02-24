package db

import (
	"fmt"

	"github.com/sunger/mygopkg/log"
	"github.com/sunger/mygopkg/model"
	"go.uber.org/zap"
)

// 键值配置分类表
type SettingCate struct {

	//模块代码
	ModuleCode string `gorm:"column:mcode;type:varchar(50);" json:"ModuleCode"`
	//名称
	Name string `gorm:"column:name;type:varchar(50);" json:"Name"`
	//显示排序
	Idx int `gorm:"column:idx;" json:"Idx"`
	model.BsModel
}

// 键值配置分类表
type SettingCateInstall struct {
	SettingCate
	Items []Settings
}

func (SettingCate) TableName() string {
	return "s_settingcate"
}

func (r *SettingCate) List() (results []SettingCate, err error) {

	err = Db.Find(&results).Error
	return results, err
}

//删除配置分类，先删除配置项
func (apps *SettingCate) Del(id string) (err error) {
	Db.Where("cate = ?", id).Delete(Settings{})
	return Db.Where("id = ?", id).Delete(SettingCate{}).Error
}
func (u *SettingCate) Insert() (id string) {
	if u.Id == "" || len(u.Id) == 0 {
		u.CreateId()
	}
	Db.Create(&u)
	return u.Id
}

// 键值配置表
type Settings struct {
	//分类id
	Cate string `gorm:"column:cate;type:varchar(50);" json:"Cate"`
	//控件json
	Ctrl string `gorm:"column:ctrl;type:varchar(500);" json:"Ctrl"`
	//Value
	Val string `gorm:"column:v;type:varchar(250);" json:"Val"`
	//显示排序
	Idx int `gorm:"column:idx;" json:"Idx"`
	model.BModel
}

func (Settings) TableName() string {
	return "s_settings"
}

func (u *Settings) Insert() (id string) {
	if u.Id == "" || len(u.Id) == 0 {
		u.CreateId()
	}
	Db.Create(&u)
	return u.Id
}

func (apps *Settings) Del(id string) (err error) {
	return Db.Where("id = ?", id).Delete(Settings{}).Error
}

// 更新实体,不存在就插入
func (apps *Settings) Update() (err error) {
	a := Settings{}
	//a.Id = apps.Id
	count := Db.Where("id=?", apps.Id).Find(&a).RowsAffected
	if count == 0 {
		log.GetLog().Info("不存在" + apps.Id)
		Db.Create(&apps)
		return err
	} else {
		log.GetLog().Info("存在" + apps.Id)
		return Db.Model(&a).Updates(map[string]interface{}{"v": apps.Val}).Error
	}
}

//根据主键获取实体
func (apps *Settings) Get(id string) (Settings, error) {
	a := Settings{}
	a.Id = id
	err := Db.Where(&a).Find(&a).Error
	return a, err
}
func (r *Settings) List(code string) (results []Settings, err error) {
	err = Db.Where("cate=?", code).Find(&results).Error
	return results, err
}

//分页方法
func (b *Settings) PageList(page, size int, filter string, sort string) ([]Settings, int) {
	db := Db

	table := b.TableName()

	if size == 0 {
		size = 20
	}

	var offset int
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * size
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

	results := make([]Settings, 0)

	rows, err := db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		log.GetLog().Error("会员分页错误", zap.String("", err.Error()))
	}

	for rows.Next() {
		var a Settings
		db.ScanRows(rows, &a)
		results = append(results, a)
	}

	var total model.PageTotal
	db.Raw(countSql).Scan(&total)
	count := total.Total

	if err != nil {
		fmt.Println(err)
	}

	return results, count
}
