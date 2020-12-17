package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/sunger/mygopkg/log"
	"go.uber.org/zap"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string, basedir string) {
	var err error
	v := viper.New()
	path := ""
	if basedir == "" {
		path, err = os.Executable()
		if err != nil {
			log.GetLog().Error(err.Error())
		}
		dir := filepath.Dir(path)
		fmt.Println("默认目录：", filepath.Join(dir, "config"))
		v.AddConfigPath(filepath.Join(dir, "config"))
	} else {
		path = basedir
		v.AddConfigPath(filepath.Join(path, "config"))
	}

	v.SetConfigType("yaml")
	v.SetConfigName(env)
	//v.AddConfigPath("../config/")
	//v.AddConfigPath("config/")
	err = v.ReadInConfig()
	if err != nil {
		log.GetLog().Fatal("解析配置文件错误", zap.String("path", path))
	}
	config = v
}

func relativePath(basedir string, path *string) {
	p := *path
	if p != "" && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}
