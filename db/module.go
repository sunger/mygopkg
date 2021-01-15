package db

import (
	//"context"
	"fmt"
	"time"
	"github.com/sunger/mygopkg/model"
	"github.com/sunger/mygopkg/log"
	"go.uber.org/zap"
	//"google.golang.org/grpc"
)

// 模块表
type Module struct {
	//Mid     string    `gorm:"column:mid;type:varchar(50);"`     //模块id
	Name          string    `gorm:"column:name;type:varchar(150);"`   //名称
	CateTp        int       `gorm:"column:catetp;type:int(1);"`       //1网站模块，2  api模块
	Domain        string    `gorm:"column:domain;type:varchar(250);"` //域名  网站模块时的值
	IsHttps       int       `gorm:"column:ishttps;type:int(1);""`     //是否https 网站模块时的值
	Tp            int       `gorm:"column:tp;type:int(1);"`           //1普通模块，2定制模块
	Path          string    `gorm:"column:path;type:varchar(50);"`    //路径，配置到caddy中，相对路径 route /mem/* {。。。 这里的mem
	Ui            string    `gorm:"column:ui;type:varchar(150);"`     //ui名称，如果前端是ng，这里可能是一个模块js名称
	No            string    `gorm:"column:no;type:varchar(100);"`     //编号
	Vs            string    `gorm:"column:vs;type:varchar(20);"`      //最新版本
	Remark        string    `gorm:"column:remark;size:250;"`          //说明
	Url           string    `gorm:"column:url;size:250;"`             //地址
	Pubdate       time.Time `gorm:"column:pub"`                       //最新发布时间
	Duedate       time.Time `gorm:"column:due"`                       //到期时间
	Img           string    `gorm:"column:img;type:varchar(150);"`    //图片
	Status        byte      `gorm:"column:status"`                    //状态,1:运行，2：停止
	InstallStatus byte      `gorm:"column:inststat"`                  //安装状态,1:运行，2：停止
	Pid           string    `gorm:"column:pid;type:varchar(50);"`     //进程id，关闭进程的时候用
	Price         float32   `gorm:"column:price"`                     //支付价格
	MemPrice      float32   `gorm:"column:memprice"`                  //会员价格
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

//插入
func (u *Module) Insert() (id string) {
	u.CreateId()
	u.Pubdate = time.Now()
	Db.Create(&u)
	return u.Id
}

// 更新实体
func (apps *Module) Update(ml Module) (err error) {
	a := Module{}
	a.Id = ml.Id
	return Db.Model(&a).Updates(map[string]interface{}{
		"name":     ml.Name,
		"catetp":   ml.CateTp,
		"tp":       ml.Tp,
		"ishttps":  ml.IsHttps,
		"domain":   ml.Domain,
		"price":    ml.Price,
		"remark":   ml.Remark,
		"memprice": ml.MemPrice,
		"vs":       ml.Vs,
		"url":      ml.Url,
		"img":      ml.Img,
		"no":       ml.No}).Error

}

//分页方法
func (b *Module) PageList(page, size int, filter string, sort string) ([]Module, int) {
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

	results := make([]Module, 0)

	rows, err := db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		log.GetLog().Error("会员等级分页错误", zap.String("", err.Error()))
	}

	for rows.Next() {
		var a Module
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
