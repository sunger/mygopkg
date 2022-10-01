package model

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"

	//"errors"
	"strconv"
	"strings"
	"time"
	//"github.com/gin-gonic/gin"
	//"github.com/sunger/mygopkg/framework/gin_"
)

//func CreateResponse() gin_.DecodeResponseFunc {
//	return func(context *gin.Context, res interface{}) error {
//		context.JSON(200, res)
//		return nil
//	}
//}
//
//func CreateQueryIdRequest() gin_.EncodeRequestFunc {
//	return func(c *gin.Context) (i interface{}, e error) {
//		bReq := &IdRequest{}
//		// err := c.ShouldBindUri(bReq)
//		bReq.Id = c.Request.FormValue("id")
//		if bReq.Id == "" {
//			return nil, errors.New("未提供query参数id")
//		}
//		return bReq, nil
//	}
//}
//
//func CreateIdRequest() gin_.EncodeRequestFunc {
//	return func(c *gin.Context) (i interface{}, e error) {
//		bReq := &IdRequest{}
//		// err := c.ShouldBindUri(bReq)
//		bReq.Id = c.Param("id")
//		if bReq.Id == "" {
//			return nil, errors.New("未提供path参数id")
//		}
//		return bReq, nil
//	}
//}
// id集合
type Ids struct {
	Ids     []string  `json:"Ids"`  //Ids
}
type CmdArg struct {
	// 模块版本
	Vs string
	// 菜单文件
	MenuFile string
	// 是否外部文件
	IsExter bool
}

type NsqModuleArg struct {
	CmdArg
	Module string //模块名称，比如mem，sys，cms等
	Cmd    string //命令，比如reload，install，remove
	Port   string //端口
}

type IdPath struct {
	// id
	Id string `uri:"id" binding:"required,gt=0,lt=30"`
}

type IdQuery struct {
	// id
	Id string `form:"id" binding:"required,gt=0,lt=30"`
}

type TreeQuery struct {
	// id
	Pid   string `form:"pid"`
	IsAll bool   `form:"isall" binding:"required"`
}

type CommResponse struct {
	// 代码
	Code int `json:"status"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg string `json:"msg"`
}

type Filter struct {
	Code string `json:"Code" binding:"required,gt=0,lt=30"`
	Tj   string `json:"Tj" binding:"required,gt=0,lt=3"`
	Val  string `json:"Val" binding:"lt=50"`
	Tp   string `json:"Tp" binding:"required,gt=0,lt=3"`
}

type Filters struct {
	Andor string   `json:"Andor" binding:"required,gt=0,lt=50"`
	Items []Filter `json:"Items"`
}

type Sorts struct {
	Code string `json:"Code" binding:"required,gt=0,lt=30"`
	Val  string `json:"Val" binding:"lt=50"`
}

/*
分页基类,每个分页基本都要这些字段
*/
type PageParams struct {
	Page int       `json:"Page" binding:"required"`
	Size int       `json:"Size" binding:"required,gt=0,lt=1000000"`
	Sort []Sorts   `json:"Sort"`
	Fts  []Filters `json:"Fts"`
}

type PageTotal struct {
	Total int64 `param:"<in:query><desc:总记录条数>"`
}

//
///*
//主键基类,每个
//*/
//type IdBase struct {
//	Id string `param:"<in:query> <required> <len: 1:50> <desc:Id (1~50 个字符)>"`
//}

//分页返回格式
type PageReturnValue struct {
	Count int64       `json:"total"`
	List  interface{} `json:"items"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type EditParam struct {
	Id    string `json:"Id" binding:"required,gt=0,lt=50"`
	Name  string `json:"Name" binding:"required,gt=0,lt=30"`
	Value string `json:"Value" binding:"required,gt=0,lt=50"`
}

func FilterItems(Items []Filter) (strs []string) {
	strs = make([]string, len(Items))
	for k, v := range Items {

		if len(v.Val) == 0 {
			continue
		}

		if v.Tp == "0" || v.Tp == "1" { // string
			if v.Tj == "0" { // bh
				strs[k] = "(" + v.Code + " = '" + v.Val + "')"
			} else if v.Tj == "1" { // baohan
				strs[k] = "(" + v.Code + " like '%" + v.Val + "%')"
			} else if v.Tj == "2" { // start
				strs[k] = "(" + v.Code + " like '" + v.Val + "%')"
			} else if v.Tj == "3" { //end
				strs[k] = "(" + v.Code + " like '%" + v.Val + "')"
			} else if v.Tj == "4" { //NBh
				strs[k] = "(" + v.Code + " not like '%" + v.Val + "')"
			} else if v.Tj == "5" { //NStart
				strs[k] = "(" + v.Code + " not like '%" + v.Val + "')"
			} else if v.Tj == "6" { //NEnd
				strs[k] = "(" + v.Code + " not like '%" + v.Val + "')"
			} else if v.Tj == "12" { //in
				strs[k] = "(" + v.Code + " in ('" + strings.Replace(v.Val, ",", "','", -1) + "'))"
			}

		} else if v.Tp == "4" { //datetime

			if v.Tj == "0" {
				strs[k] = "(" + v.Code + " = '" + v.Val + "')"
			} else if v.Tj == "7" { //Lt
				strs[k] = "(" + v.Code + " < '" + v.Val + "')"
			} else if v.Tj == "8" { //Lte
				strs[k] = "(" + v.Code + " <= '" + v.Val + "')"
			} else if v.Tj == "9" { //Gt
				strs[k] = "(" + v.Code + " > '" + v.Val + "')"
			} else if v.Tj == "10" { //Gte
				strs[k] = "(" + v.Code + " >= '" + v.Val + "')"
			}
		} else if v.Tp == "5" { //位运算
			strs[k] = "((" + v.Code + " & " + v.Val + ") > 0)"
		} else { //bool number
			if v.Tj == "0" {
				strs[k] = "(" + v.Code + " = " + v.Val + ")"
			} else if v.Tj == "7" { //Lt
				strs[k] = "(" + v.Code + " < " + v.Val + ")"
			} else if v.Tj == "8" { //Lte
				strs[k] = "(" + v.Code + " <= " + v.Val + ")"
			} else if v.Tj == "9" { //Gt
				strs[k] = "(" + v.Code + " > " + v.Val + ")"
			} else if v.Tj == "10" { //Gte
				strs[k] = "(" + v.Code + " >= " + v.Val + ")"
			} else if v.Tj == "12" { //in
				strs[k] = "(" + v.Code + " in (" + v.Val + "))"
			} else if v.Tj == "13" { //多少天之内的日期查询
				days, _ := strconv.Atoi(v.Val)
				if days > 0 {
					now := time.Now()
					fmt := "-" + strconv.Itoa(days*24) + "h"
					// d, _ := time.ParseDuration("-24h")
					d, _ := time.ParseDuration(fmt)
					d1 := now.Add(d)

					start_ := d1.Format("2006-01-02 15:04:05")
					end_ := now.Format("2006-01-02 15:04:05")

					strs[k] = "(" + v.Code + " Between  '" + start_ + "' and '" + end_ + "')"
				}

			}
		}

	}
	return strs
}

func FilterStr(Items []Filters) string {

	ln := len(Items)
	if ln == 0 {
		return ""
	}
	filters := make([]string, ln)

	for k, v := range Items {
		// filters[k] = v.Andor + " 123"

		strs := FilterItems(v.Items)

		//去掉里面的空值
		var strsValide []string
		for _, v1 := range strs {
			if len(v1) > 0 {
				strsValide = append(strsValide, v1)
			}
		}

		ln2 := len(strsValide)
		if ln2 > 0 {
			filters[k] = strings.Join(strsValide, " "+v.Andor+" ")
		}
	}
	filterstr := strings.Join(filters, " and ")

	if len(filterstr) < 7 {
		return ""
	}
	return filterstr
}

func GetFlts(p PageParams) (strs []string) {
	if p.Page == 0 {
		p.Page = 1
	}
	lenSort := len(p.Sort)
	// filters := make([]string, len(p.Fts))

	orderstr := ""

	if lenSort > 0 {
		orders := make([]string, len(p.Sort))
		//排序字段
		for k, v := range p.Sort {

			orders[k] = v.Code + " " + v.Val
		}

		strings.Join(orders, ",")
	}

	//查询字段
	// for k, v := range p.Fts {
	// 	// filters[k] = v.Andor + " 123"
	// 	filters[k] = strings.Join(FilterItems(v.Items), " "+v.Andor+" ")

	// 	beego.Debug(filters[k])
	// }
	strs = make([]string, 2)

	filterstr := FilterStr(p.Fts)
	/*
	   HasPrefix 判断字符串 s 是否以 prefix 开头：
	   strings.HasPrefix(s, prefix string) bool
	   HasSuffix 判断字符串 s 是否以 suffix 结尾：
	   strings.HasSuffix(s, suffix string) bool
	*/
	filterstr = strings.TrimSpace(filterstr)
	if len(filterstr) > 0 {

		if strings.HasPrefix(filterstr, "or") || strings.HasPrefix(filterstr, "and") {
			strs[0] = "1=1 " + filterstr
		} else {
			strs[0] = filterstr
		}
		// if strings.HasPrefix(filterstr, "and") {
		// 	strs[0] = "1=1 " + filterstr
		// }

	} else {
		strs[0] = filterstr
	}

	strs[1] = orderstr

	return strs
}

/////////////////////////////////////////////
//https://www.jianshu.com/p/2d841ffae6af
///自定义日期类型

const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

/*
yann：// 指定解析的格式，这样写会出现时区问题
now, err := time.Parse("+TimeFormat+", string(data))
// 最好指定一下时区
loc, _ := time.LoadLocation("Asia/Shanghai")
now, err := time.ParseInLocation("+TimeFormat+", string(data), loc)

*/
type BModel2 struct {
	gorm.Model
	Id string `gorm:"column:id;primary_key;type:varchar(50)" json:"Id"` //主键
}

//通用分页方法，只支持一个表或者视图的查询
func PageList(db *gorm.DB, table string, page, size int, filter string, sort string) ([]map[string]interface{}, int64) {

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

	arrSql := fmt.Sprintf("SELECT * FROM %s WHERE %s order by %s limit %d offset %d",
		table, filter, sort, size, offset)

	countSql := fmt.Sprintf("SELECT count(0) as total FROM %s WHERE %s",
		table, filter)

	results := make([]map[string]interface{}, 0)

	rows, err := db.Raw(arrSql).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var a map[string]interface{}
		db.ScanRows(rows, &a)
		results = append(results, a)
	}

	var total PageTotal
	db.Raw(countSql).Scan(&total)
	count := total.Total

	if err != nil {
		fmt.Println(err)
	}

	return results, count
}
