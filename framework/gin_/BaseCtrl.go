package gin_

import (
	"github.com/gin-gonic/gin"
	"github.com/sunger/mygopkg/log"
)

//防止重复解析，缓存token：用户信息
var logins = make(map[string]UserInfo,0)
// 控制器基类
type BaseCtrl struct{}

//获取当前登录的用户信息
func (m *BaseCtrl) User(c *gin.Context) UserInfo {

	token := c.Request.Header.Get("token")
	if token == "" {
		log.GetLog().Error("token为空")
		return UserInfo{}
	}

	if v, ok := logins[token]; ok {
		return v
	} else {
		user, err := ParseToken(token)
		if err != nil {
			log.GetLog().Error("解析token错误")
			return UserInfo{}
		}
		//保存到缓存中
		logins[token] = user
		return user
	}
}

//清空
func (m *BaseCtrl) Clear()  {
	logins = make(map[string]UserInfo,0)
}

//重置
func (m *BaseCtrl) Reset(oldtoken, newtoken string, user UserInfo)  {
	if _, ok := logins[oldtoken]; ok {
		delete(logins, oldtoken)
	}
	logins[newtoken] = user
}

