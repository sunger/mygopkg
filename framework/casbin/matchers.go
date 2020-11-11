package casbin

import (
	"fmt"
	"strings"
)

func init() {

	E.AddFunction("actionMatch", func(arguments ...interface{}) (i interface{}, e error) {
		if len(arguments) == 2 {
			k1, k2 := arguments[0].(string), arguments[1].(int)
			fmt.Println("参数：", k1, k2)
			return ActionMatch(k1, k2), nil
		}
		fmt.Println("请求方法不匹配::")
		return nil, fmt.Errorf("请求方法不匹配")
	})

}

// 自定义匹配函数 匹配web请求方法
// auth：数据库中配置的action权限集合,数字类型的
// act: web中传入的函数
func ActionMatch(act string, auth int) bool {
	v := Actions[strings.ToLower(act)]
	return (v & auth) > 0
}
