package gin_

import (
	"github.com/gin-gonic/gin"
	"github.com/sunger/mygopkg/log"
)

// 控制器基类
type BaseCtrl struct{}

//获取当前登录的用户信息
func (m *BaseCtrl) User(c *gin.Context) UserInfo {

	token := c.Request.Header.Get("token")
	if token == "" {
		log.GetLog().Error("token为空")
		return UserInfo{}
	}
	user, err := ParseToken(token)
	if err != nil {
		log.GetLog().Error("解析token错误")
		return UserInfo{}
	}
	return user
}


