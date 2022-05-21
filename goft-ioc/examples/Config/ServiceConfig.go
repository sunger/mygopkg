package Config

import (
	//"github.com/sunger/goft-ioc/examples/services"
	"github.com/sunger/mygopkg/goft-ioc/examples/services"
)

type ServiceConfig struct {
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}
func (this *ServiceConfig) OrderService() *services.OrderService {
	return services.NewOrderService()
}
