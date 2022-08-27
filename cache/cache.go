package cache

import (
	"time"
	"fmt"
	cc "github.com/patrickmn/go-cache"
)
var usercache_ *cc.Cache  	//登录用户缓存
var rolecache_ *cc.Cache 	//角色缓存
var envcache_ *cc.Cache		//全局配置变量缓存

//用户相关缓存
type UserCache struct {
	IsSuper bool //是否超级管理员
	RoleIds []string //角色id相关缓存
	Ext map[string]interface{} //其他相关缓存
}
//缓存用户信息
func (m *UserCache) Set(id string, issuper bool,roleIds []string,ext map[string]interface{}) {
	//先查询
	rc, exist := m.Get(id)
	if exist {
		rc.IsSuper = issuper
		rc.RoleIds = roleIds
		rc.Ext = ext
	}else {
		userCache:= &  UserCache{}
		userCache.IsSuper = issuper
		userCache.RoleIds = roleIds
		userCache.Ext = ext
		usercache_.Set(id, userCache, cc.DefaultExpiration)
	}

}
//缓存用户信息额外信息
func (m *UserCache) SetExt(id, key string, v interface{}) {
	//先查询
	rc, exist := m.Get(id)
	if exist {
		rc.Ext[key] = v
	}else {
		userCache:= &  UserCache{}
		userCache.IsSuper = false
		userCache.RoleIds = make([]string,0)
		userCache.Ext = map[string]interface{}{
			key: v,
		}
		usercache_.Set(id, userCache, cc.DefaultExpiration)
	}

}

//获取用户信息
func (m *UserCache) Get(id string) (UserCache,bool) {
	i, b := usercache_.Get(id)
	if b {
		uc, ok := (i).(UserCache)
		if ok {
			return uc,b
		} else {
			fmt.Println("UserCache缓存数据转换失败"+id)
		}
	}
	return *m, b
}
//删除用户信息
func (m *UserCache) Delete(id string) bool {
	_, b := usercache_.Get(id)
	if b {
		usercache_.Delete(id)
	}
	return  b
}


//角色相关缓存
type RoleCache struct {
	Id string //角色id
	IsSuper bool //是否超级管理员角色
	Ext map[string]interface{} //其他相关缓存
}

//缓存角色信息
func (m *RoleCache) Set(id string,issuper bool,ext map[string]interface{}) {
	//先查询
	rc, exist := m.Get(id)
	if exist {
		rc.IsSuper = issuper
		rc.Ext = ext
	}else{
		roleCache:= &  RoleCache{}
		roleCache.IsSuper = issuper
		roleCache.Ext = ext
		rolecache_.Set(id, roleCache, cc.DefaultExpiration)
	}

}
//缓存角色信息额外信息
func (m *RoleCache) SetExt(id, key string, v interface{}) {
	//先查询
	rc, exist := m.Get(id)
	if exist {
		rc.Ext[key] = v
	}else {
		roleCache:= &  RoleCache{}
		roleCache.IsSuper = false
		roleCache.Ext = map[string]interface{}{
			key: v,
		}
		usercache_.Set(id, roleCache, cc.DefaultExpiration)
	}
}

//获取角色信息
func (m *RoleCache) Get(id string) (RoleCache,bool) {
	i, b := rolecache_.Get(id)
	if b {
		uc, ok := (i).(RoleCache)
		if ok {
			return uc,b
		} else {
			fmt.Println("RoleCache缓存数据转换失败:"+id)
		}
	}
	return *m, b
}
//删除用户信息
func (m *RoleCache) Delete(id string) bool {
	_, b := rolecache_.Get(id)
	if b {
		rolecache_.Delete(id)
	}
	return  b
}


//全局变量相关缓存
type EnvCache struct {
	Key string //Key
}
func (m *EnvCache) Get(id string) (interface{},bool) {
	return envcache_.Get(id)
}
//缓存全局变量信息
func (m *EnvCache) Set(key string, v interface{}) {
	envcache_.Set(key, v, cc.DefaultExpiration)
}

//删除用户信息
func (m *EnvCache) Delete(id string) bool {
	_, b := envcache_.Get(id)
	if b {
		envcache_.Delete(id)
	}
	return  b
}


var RqRoleCache RoleCache
var RqUserCache UserCache
var RqEnvCache EnvCache

//初始化缓存
func init() {
	//默认10分钟缓存
	usercache_ = cc.New(30*time.Minute, 30*time.Minute)
	rolecache_ = cc.New(60*time.Minute, 60*time.Minute)
	envcache_ = cc.New(120*time.Minute, 120*time.Minute)
}
