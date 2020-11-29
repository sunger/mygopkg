package tools

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/sunger/mygopkg/config"
)

func TestBsModel(t *testing.T) {

	imageUtil := ImageUtil{}
	imageUtil.GetWH("")
}

func TestUpload(t *testing.T) {
	env := flag.String("e", "development", "用法：xxx -e development|product|test")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	// 初始化配置文件,这里非常重要，默认是development，理论上这里可以是任何值，只要存在类似：development。yaml文件
	// 在微服务里可以针对每个模块启动对应的数据库连接，也就是不同的数据库
	// 在本地单击部署支持不同的数据库，在云端可以考虑把这些配置文件存配置中心

	config.Init(*env)
	c := config.GetConfig()
	defaultoss := c.GetString("oss.default")
	fmt.Println(defaultoss)
}
