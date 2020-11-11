package cfg

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

// Yaml struct of yaml
type YamlCfg struct {
	Gin struct {
		Domain string `yaml:"domain"`
		Port   string `yaml:"port"`
	}
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
	Sqlite struct {
		Name string `yaml:"name"`
	}
	Cache struct {
		Enable bool     `yaml:"enable"`
		List   []string `yaml:"list,flow"`
	}
}

// 构造函数读 配置文件
func NewYamlCfg() *YamlCfg {
	conf := &YamlCfg{}
	yamlFile, err := ioutil.ReadFile("cfg.yaml")

	// conf := new(module.Yaml1)
	// yamlFile, err := ioutil.ReadFile("test.yaml")

	// conf := new(module.Yaml2)
	//  yamlFile, err := ioutil.ReadFile("test1.yaml")

	// log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println("conf", conf)
	// log.Println("conf", resultMap)
	return conf
}
