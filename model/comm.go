package model

type IdRequest struct {
	// id
	Id string `uri:"id" binding:"required"`
}

type CommResponse struct {
	// 代码
	Code int `json:"code"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg string `json:"msg"`
}

/*
分页基类,每个分页基本都要这些字段
*/
type ListBase struct {
	Page  int    `param:"<in:query><desc:当前页>"`
	Size  int    `param:"<in:query><desc:每页记录数>"`
	Sort  string `param:"<in:query><desc:排序字段>"`
	Order int    `param:"<in:query><desc:排序类型:1:asc,0:desc>"`
}

type Filter struct {
	Code string `param:"<in:query><desc:字段名称>"`
	Tj   string `param:"<in:query><desc:条件（>,<,=）等>"`
	Val  string `param:"<in:query><desc:字段值>"`
	Tp   string `param:"<in:query><desc:字段数据类型>"`
}

type Filters struct {
	Andor string   `param:"<in:query><desc:and,or>"`
	Items []Filter `param:"<in:query><desc:条件项数组>"`
}

type Sorts struct {
	Code string `param:"<in:query><desc:字段名称>"`
	Val  string `param:"<in:query><desc:字段值>"`
}

/*
分页基类,每个分页基本都要这些字段
*/
type PageParams struct {
	Page int       `param:"<in:query><desc:当前页>"`
	Size int       `param:"<in:query><desc:每页记录数>"`
	Sort []Sorts   `param:"<in:query><desc:排序字段集合>"`
	Fts  []Filters `param:"<in:query><desc:搜索条件>"`
}

type PageTotal struct {
	Total int `param:"<in:query><desc:总记录条数>"`
}

/*
主键基类,每个
*/
type IdBase struct {
	Id string `param:"<in:query> <required> <len: 1:50> <desc:Id (1~50 个字符)>"`
}

//分页返回格式
type PageReturnValue struct {
	Count int
	List  interface{}
}

type EditParam struct {
	Id    string `param:"<in:query><desc:更新主键>"`
	Name  string `param:"<in:query><desc:更新字段名称>"`
	Value string `param:"<in:query><desc:更新字段值>"`
}
