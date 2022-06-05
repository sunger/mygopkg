package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// 所有模型类的基类
type BModel struct {
	Id    string `gorm:"column:id;primary_key;type:varchar(50)" json:"Id"` 	//主键
	Tid  string `gorm:"column:tid;type:varchar(50);" json:"Tid"`   					//租户id
	DbKey string `gorm:"-" json:"-"`                                       	//数据库key，根据此key指向对应的数据库
}

// 所有模型类的基类,带公司和店铺
type BsModel struct {
	BModel
	Cid  string `gorm:"column:cid;type:varchar(50);" json:"Cid"`   //公司id
	Sid  string `gorm:"column:sid;type:varchar(50);" json:"Sid"`   //店铺id
	Ctid string `gorm:"column:ctid;type:varchar(50);" json:"Ctid"` //创建者id
}

// 创建id：默认使用时间戳方式生成
func (r *BModel) CreateId() {
	idtp := 0 // id生成方式

	if idtp == 0 {
		// 时间戳（秒）：1530027865; 10位数的时间戳是以 秒 为单位；
		r.Id = strconv.FormatInt(time.Now().Unix(), 10)
	} else if idtp == 1 {
		id_ := uuid.NewV4()
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
	CreatedAt time.Time  `gorm:"column:createdat" json:"CreatedAt"` //创建时间
	UpdatedAt time.Time  `gorm:"column:updatedat" json:"UpdatedAt"` //最后修改时间
	DeletedAt *time.Time `gorm:"column:deletedat" json:"DeletedAt"` //删除时间
	BsModel
}

// 所有树模型类的基类（6个分类相关的字段）
type TreeModel struct {
	Grp string `gorm:"column:grp;type:varchar(20);" json:"Grp"` //树分组标识
	Lft int    `gorm:"column:lft;type:int(11);" json:"Lft"`     //树左节点
	Lv  int    `gorm:"column:lv;type:int(2);" json:"Lv"`        //树层级
	Idx int    `gorm:"column:idx;type:int(6);" json:"Idx"`      //树层中排序
	Rgt int    `gorm:"column:rgt;type:int(11);" json:"Rgt"`     //树右节点
	Pid string `gorm:"column:pid;type:varchar(50);" json:"Pid"` //树父节点id
	Tid string `gorm:"column:tid;type:varchar(50);" json:"Tid"` //租户id
}

// 防止此代码在所有树型记录中被重复编写（添加动作中）
// 树模型类的插入,具体的插入自己操作，这里只负责修改生成树结构，此方法在自己的插入动作之前操作，且应该和自己的插入在同一个事务中
func (r *TreeModel) BeforInsertTree(tableName string, db *gorm.DB) (parent TreeModel,err error) {
	var count int64 = 0
	pnode := TreeModel{}
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
	

		// if pnode, err = r.Get(r.Pid); err != nil {
		// 	return errors.New("父节点对应的记录不存在")
		// }

		db.Table(tableName).Where("id=?", r.Pid).Select("rgt, lft, lv, grp, idx, pid, tid").First(&pnode).Count(&count)
		if count == 0 {
			return pnode , errors.New("父节点对应的记录不存在")
		}

		//查询此父节点下的最大子节点
		//如果存在，插入最多子节点右边，否则插入父节点下的作为第一个子节点
		maxSubNd := TreeModel{}

		// gorm.MustDB().Order("idx desc").Where("pid=?", r.Pid).First(&maxSubNd).Count(&count)

		db.Table(tableName).Order("idx desc").Where("pid=?", r.Pid).Select("rgt, lft, lv, grp, idx, pid, tid").First(&maxSubNd).Count(&count)

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

	return pnode, nil
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

	db.Exec("delete from "+tableName+" WHERE id = ?", id)
	db.Exec("update "+tableName+" set rgt = rgt - 1, lft = lft - 1 WHERE grp = ? and lft BETWEEN ? and ?", grp, lft, rgt)
	db.Exec("update "+tableName+" set rgt = rgt - 2 WHERE grp = ? and rgt > ?", grp, rgt)
	db.Exec("update "+tableName+" set lft = lft - 2 WHERE grp = ? and lft > ?", grp, rgt)
	return err
}

/////////////

//MyTime 自定义时间
type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = MyTime(t1)
	return err
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t MyTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *MyTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

//通用分页方法
