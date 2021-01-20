package db

import (
	"strings"
	//"context"
	"fmt"
	"time"

	"github.com/sunger/mygopkg/log"
	"github.com/sunger/mygopkg/model"
	"go.uber.org/zap"
	//"google.golang.org/grpc"
)

// 模块表
type Module struct {
	//Mid     string    `gorm:"column:mid;type:varchar(50);"`     //模块id
	Name          string    `gorm:"column:name;type:varchar(150);" json:"Name"`     //名称
	CateTp        int       `gorm:"column:catetp;type:int(1);" json:"CateTp"`       //1网站模块，2  api模块
	FeeTp         int       `gorm:"column:feetp;type:int(1);" json:"FeeTp"`         //1免费模块，2  收费模块
	Domain        string    `gorm:"column:domain;type:varchar(250);" json:"Domain"` //域名  网站模块时的值
	IsHttps       int       `gorm:"column:ishttps;type:int(1);"" json:"IsHttps"`    //是否https 网站模块时的值
	Tp            int       `gorm:"column:tp;type:int(1);" json:"Tp"`               //1普通模块，2定制模块
	Path          string    `gorm:"column:path;type:varchar(50);" json:"Path"`      //路径，配置到caddy中，相对路径 route /mem/* {。。。 这里的mem
	Ui            string    `gorm:"column:ui;type:varchar(150);" json:"Ui"`         //ui名称，如果前端是ng，这里可能是一个模块js名称
	No            string    `gorm:"column:no;type:varchar(100);" json:"No"`         //编号
	Vs            string    `gorm:"column:vs;type:varchar(20);" json:"Vs"`          //最新版本
	CVs           string    `gorm:"column:cvs;type:varchar(20);" json:"CVs"`        //当前版本
	Remark        string    `gorm:"column:remark;size:250;" json:"Remark"`          //说明
	Url           string    `gorm:"column:url;size:250;" json:"Url"`                //地址
	Pubdate       time.Time `gorm:"column:pub" json:"Pubdate"`                      //最新发布时间
	Duedate       time.Time `gorm:"column:due" json:"Duedate"`                      //到期时间
	Img           string    `gorm:"column:img;type:varchar(150);" json:"Img"`       //图片
	Idx           int       `gorm:"column:idx;type:int(6);" json:"Idx"`             //排序
	Status        byte      `gorm:"column:status" json:"Status"`                    //状态,1:运行，2：停止
	InstallStatus byte      `gorm:"column:inststat" json:"Pid"`                     //安装状态,1:已安装，2：未安装
	Pid           string    `gorm:"column:pid;type:varchar(50);" json:"Idx"`        //进程id，关闭进程的时候用
	//Price         float32   `gorm:"column:price"`                     //支付价格
	//MemPrice      float32   `gorm:"column:memprice"`                  //会员价格
	model.BaseModel
}

func (Module) TableName() string {
	return "rq_module"
}

//根据主键获取实体
func (r *Module) Get(id string) (Module, error) {
	a := Module{}
	a.Id = id
	err := Db.Where(&a).Find(&a).Error
	return a, err
}
func (u *Module) FindByPath(path string) (Module, error) {
	user := Module{}
	err := Db.Find(&user, "path = ? ", path).Error
	return user, err
}

//
func (r *Module) GetAll() (results []Module, err error) {

	err = Db.Find(&results).Error
	return results, err
}

// 停止启动模块
func (apps *Module) ChangeStatus() (err error) {
	a := Module{}
	a.Id = apps.Id
	////rpc先结束当前模块，然后删除数据库中的记录
	//conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//cmd := "stop"
	//
	//if apps.Status == 1 {
	//	cmd = "start"
	//}
	//
	//cmdClient := services.NewCmdServiceClient(conn)
	//_, err = cmdClient.ExecCmdResponse(context.Background(), &services.CmdRequest{Module: apps.Path, Cmd: cmd})
	//if err != nil {
	//	return err
	//}

	return Db.Model(&a).Updates(map[string]interface{}{
		"status": apps.Status}).Error
}

//删除
func (u *Module) Del() error {
	////rpc先结束当前模块，然后删除数据库中的记录
	//conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//cmdClient := services.NewCmdServiceClient(conn)
	//_, err = cmdClient.ExecCmdResponse(context.Background(), &services.CmdRequest{Module: u.Path, Cmd: "remove"})
	//if err != nil {
	//	return err
	//}

	return Db.Delete(&u).Error
}

func (u *Module) FindByName(name string) (Module, error) {
	user := Module{}
	err := Db.Where("name = ?", name).First(&user).Error
	return user, err
}

func (u *Module) FindByPathOrNo(path, no string) (Module, error) {
	user := Module{}
	err := Db.Find(&user, "path = ? or no = ?", path, no).Error
	return user, err
}

//插入
func (u *Module) Insert() (id string, err error) {
	//
	//err = Db.Where("path = ?", u.Path).Error
	//
	//if !errors.Is(err, gorm.ErrRecordNotFound) {
	//	return "", errors.New("路径path不能重复")
	//}
	//err = Db.Where("no = ?", u.No).Error
	//
	//if !errors.Is(err, gorm.ErrRecordNotFound) {
	//	return "", errors.New("编号no不能重复")
	//}

	u.CreateId()
	u.Pubdate = time.Now()
	u.InstallStatus = 2
	u.Status = 2
	Db.Create(&u)


	return u.Id, nil
}

// 更新实体
func (apps *Module) Update(ml Module) (err error) {
	a := Module{}
	a.Id = ml.Id
	return Db.Model(&a).Updates(map[string]interface{}{
		"name":     ml.Name,
		"catetp":   ml.CateTp,
		"tp":       ml.Tp,
		"feetp":    ml.FeeTp,
		"ishttps":  ml.IsHttps,
		"inststat": ml.InstallStatus,
		"status":   ml.Status,
		"domain":   ml.Domain,
		//"price":    ml.Price,
		"remark": ml.Remark,
		//"memprice": ml.MemPrice,
		"vs":   ml.Vs,
		"url":  ml.Url,
		"img":  ml.Img,
		"ui":   ml.Ui,
		"pub":  ml.Pubdate,
		"due":  ml.Duedate,
		"path": ml.Path,
		"no":   ml.No}).Error

}

// 更新实体
func (apps *Module) UpdateVs(ml Module) (err error) {
	a := Module{}
	a.Id = ml.Id
	return Db.Model(&a).Updates(map[string]interface{}{
		"pub": time.Now(),
		"vs":  ml.Vs}).Error

}

// 更新实体，指定字段指定值
func (r *Module) UpdateByFeild(id, name, value string) (err error) {
	r.Id = id

	return Db.Model(r).Update(strings.ToLower(name), value).Error
}

//分页方法
func (b *Module) PageList(page, size int, filter string, sort string) ([]Module, int) {

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

	results := make([]Module, 0)

	rows, err := Db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		log.GetLog().Error("模块分页错误", zap.String("", err.Error()))
	}

	for rows.Next() {
		var a Module
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

//列表方法
func (b *Module) List(filter string, sort string) []Module {

	table := b.TableName()

	if len(filter) == 0 {
		filter = "1=1"
	}
	if len(sort) == 0 {
		sort = "idx"
	}

	arrSql := fmt.Sprintf("SELECT * FROM %s WHERE %s order by %s",
		table, filter, sort)

	results := make([]Module, 0)

	rows, err := Db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		log.GetLog().Error("模块列表错误", zap.String("", err.Error()))
	}

	for rows.Next() {
		var a Module
		Db.ScanRows(rows, &a)
		results = append(results, a)
	}
	if err != nil {
		fmt.Println(err)
	}

	return results
}
