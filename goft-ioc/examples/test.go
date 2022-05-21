package main

import (
	"fmt"
	Injector "github.com/sunger/mygopkg/goft-ioc"
	"github.com/sunger/mygopkg/goft-ioc/examples/Config"
	"github.com/sunger/mygopkg/goft-ioc/examples/services"

	//"github.com/sunger/goft-ioc"
	//"github.com/sunger/goft-ioc/examples/Config"
	//"github.com/sunger/goft-ioc/examples/services"
)

func main() {
	serviceConfig := Config.NewServiceConfig()

	Injector.BeanFactory.Config(serviceConfig) //展开方法
	//  BeanFactory.Set()
	{
		//这里 测试 userServices
		userService := services.NewUserService()
		Injector.BeanFactory.Apply(userService) //处理依赖
		fmt.Println(userService.Order.Name())
		userService.GetUserInfo(3)

	}
}
