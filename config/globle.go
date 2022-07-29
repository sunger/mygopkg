package config

/*
全局变量
*/
type Global struct {
	DbId     string //数据库链接字符串id
	DbUserTp int    //数据库链接字符串使用方式：1:默认数据库（命令行参数录入或者程序中默认值）,2:用户登录值
}

var Glb Global

func init() {
	Glb.DbId = "1"
	Glb.DbUserTp = 1
}
