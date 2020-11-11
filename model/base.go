package model

import (
	"errors"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// 所有模型类的基类
type BsModel struct {
	Id  string `gorm:"column:id;primary_key;type:varchar(50)"` //主键
	Cid string `gorm:"column:cid;type:varchar(50);"`           //公司id
	Sid string `gorm:"column:sid;type:varchar(50);"`           //店铺id
}

// 创建id：默认使用时间戳方式生成
func (r *BsModel) CreateId() {
	idtp := 0 // id生成方式

	if idtp == 0 {
		// 时间戳（秒）：1530027865; 10位数的时间戳是以 秒 为单位；
		r.Id = strconv.FormatInt(time.Now().Unix(), 10)
	} else if idtp == 1 {
		id_, _ := uuid.NewV4()
		r.Id = id_.String()

		r.Id = strings.Replace(r.Id, "-", "", -1)

	} else if idtp == 2 {
		// 时间戳（纳秒）：1530027865231834600; 19位数的时间戳是以 纳秒 为单位；
		r.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	} else if idtp == 3 {
		// 时间戳（毫秒）：1530027865231;  13位数的时间戳是以 毫秒 为单位；
		r.Id = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	}

}

// 所有模型类的基类（多3个时间字段）
type BaseModel struct {
	CreatedAt time.Time  `gorm:"column:createdat"` //创建时间
	UpdatedAt time.Time  `gorm:"column:updatedat"` //最后修改时间
	DeletedAt *time.Time `gorm:"column:deletedat"` //删除时间
	BsModel
}

// 所有树模型类的基类（6个分类相关的字段）
type TreeModel struct {
	Grp string `gorm:"column:grp;size:20"`           //树分组标识
	Lft int    `gorm:"column:lft;type:int(11);"`     //树左节点
	Lv  int    `gorm:"column:lv;type:int(2);"`       //树层级
	Idx int    `gorm:"column:idx;type:int(6);"`      //树层中排序
	Rgt int    `gorm:"column:rgt;type:int(11);"`     //树右节点
	Pid string `gorm:"column:pid;type:varchar(50);"` //树父节点id
}

// 防止此代码在所有树型记录中被重复编写（添加动作中）
// 树模型类的插入,具体的插入自己操作，这里只负责修改生成树结构，此方法在自己的插入动作之前操作，且应该和自己的插入在同一个事务中
func (r *TreeModel) BeforInsertTree(tableName string, db *gorm.DB) (err error) {
	var count int64 = 0
	if len(r.Pid) == 0 {

		// gorm.MustDB().Model(&r).Count(&count)

		// if count > 0 {
		// 	return errors.New("已经存在记录，父节点不可为空")
		// } else {
		//跟节点
		r.Lft = 1
		r.Rgt = 2
		r.Lv = 1
		r.Idx = 1
		// }

	} else {
		pnode := TreeModel{}

		// if pnode, err = r.Get(r.Pid); err != nil {
		// 	return errors.New("父节点对应的记录不存在")
		// }

		db.Table(tableName).Where("id=?", r.Pid).Select("rgt, lft, lv, grp, idx, pid").First(&pnode).Count(&count)
		if count == 0 {
			return errors.New("父节点对应的记录不存在")
		}

		//查询此父节点下的最大子节点
		//如果存在，插入最多子节点右边，否则插入父节点下的作为第一个子节点
		maxSubNd := TreeModel{}

		// gorm.MustDB().Order("idx desc").Where("pid=?", r.Pid).First(&maxSubNd).Count(&count)

		db.Table(tableName).Order("idx desc").Where("pid=?", r.Pid).Select("rgt, lft, lv, grp, idx, pid").First(&maxSubNd).Count(&count)

		if count == 0 {
			// beego.Info("没有记录")
			r.Lv = pnode.Lv + 1
			r.Idx = 1
			r.Lft = pnode.Lft + 1
			r.Rgt = pnode.Lft + 2
			db.Exec("update "+tableName+" set rgt = rgt + 2 WHERE rgt > ?;", pnode.Lft)
			db.Exec("update "+tableName+" set lft = lft + 2 WHERE lft > ?;", pnode.Lft)
		} else {
			// beego.Info("有记录")
			db.Exec("update "+tableName+" set rgt = rgt + 2 WHERE rgt > ?;", maxSubNd.Rgt)
			db.Exec("update "+tableName+" set lft = lft + 2 WHERE lft > ?;", maxSubNd.Rgt)
			r.Lft = maxSubNd.Rgt + 1
			r.Rgt = maxSubNd.Rgt + 2
			r.Lv = maxSubNd.Lv
			r.Idx = maxSubNd.Idx + 1
		}

		// 如果是子节点，继承父节点的Group
		r.Grp = pnode.Grp
	}

	return nil
}

//删除节点和子节点
func (r *TreeModel) DelAll(tableName, grp, id string, rgt, lft int, db *gorm.DB) (err error) {

	myWidth := rgt - lft + 1
	db.Exec("delete from "+tableName+" WHERE grp = ? and lft  BETWEEN ? AND  ?", grp, lft, rgt)
	db.Exec("update "+tableName+" set rgt = rgt - ? WHERE grp = ? and rgt > ?", myWidth, grp, rgt)
	db.Exec("update "+tableName+" set lft = lft - ? WHERE grp = ? and lft > ?", myWidth, grp, rgt)
	return err
}

//删除当前节点
func (r *TreeModel) DelOne(tableName, grp, id string, rgt, lft int, db *gorm.DB) (err error) {

	db.Exec("delete from "+tableName+" WHERE grp = ? lft = ?", rgt, lft)
	db.Exec("update "+tableName+" set rgt = rgt - 1，lft = lft - 1 WHERE grp = ? and lft BETWEEN ? and ?", grp, lft, rgt)
	db.Exec("update "+tableName+" set rgt = rgt - 2 WHERE grp = ? and rgt > ?", grp, rgt)
	db.Exec("update "+tableName+" set lft = lft - 2 WHERE grp = ? and lft > ?", grp, rgt)
	return err
}
