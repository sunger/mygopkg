package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/sunger/mygopkg/log"
	"go.uber.org/zap"
)

var config *viper.Viper


func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	//fmt.Println( l)
	return string(runes[pos:l])
}

//上级目录
func ParentDir(dirctory string) string {
	//fmt.Println("当前目录：" + dirctory)
	d := filepath.ToSlash(dirctory)
	fmt.Println("转换后的当前目录：" + d)
	return substr(d, 0, strings.LastIndex(d, "/"))
}

//当前目录
func Dir() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}

//下级目录
func SubDir(dir, sub string) string {
	path := filepath.Join(dir, sub)
	return path
}

//路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

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

	cfgFullPath := filepath.Join(filepath.Join(path, "config"), env)
	cfgFullPath = cfgFullPath + ".yaml"
	exist, _ := PathExists(cfgFullPath)
	//fmt.Println(exist)
	if exist {
		fmt.Println("yaml配置文件：", cfgFullPath)
	} else {
		fmt.Println("yaml配置文件不存在：", cfgFullPath)
	}

	v.SetConfigType("yaml")
	v.SetConfigName(env)
	//v.AddConfigPath("../config/")
	//v.AddConfigPath("config/")
	err = v.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
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
	if config == nil {
		fmt.Println("请先config.Init初始化配置文件")
		return nil
	}
	return config
}
