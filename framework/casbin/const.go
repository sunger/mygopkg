package casbin

import (
	"github.com/casbin/casbin/v2"
)

//全局变量
var E *casbin.Enforcer

//web请求方法枚举
var Actions = map[string]int{
	"get":    1,  //get方法
	"post":   2,  //post方法
	"put":    4,  //put方法
	"delete": 8,  //delete方法
	"option": 16, //option方法
}
