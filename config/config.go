package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string `yaml:"name"`
	Env    string `yaml:"env"`
	Server Server `yaml:"server"`
	Redis  Redis  `yaml:"redis"`
	Log    Log    `yaml:"log"`
}
type Server struct {
	Port string `yaml:"port"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}
type Log struct {
	Filename string `yaml:"filename"`
}

// 配置信息缓存
var Cfg Config

// 初始化 config 配置
func Init() {

	// 解析默认基础配置文件
	filename := filepath.Join("config", "config.yml")
	yml, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err1 := yaml.Unmarshal(yml, &Cfg)
	if err1 != nil {
		panic(err1)
	}
	// fmt.Println(filename, Cfg)

	// 原始 env名称
	env := Cfg.Env
	switch Cfg.Env {
	case "dev":
		Cfg.Env = gin.DebugMode
	case "test":
		Cfg.Env = gin.TestMode
	case "prod":
		Cfg.Env = gin.ReleaseMode
	default:
		Cfg.Env = gin.DebugMode
		fmt.Println("error, env is set to debug mode")
	}
	if env != "dev" {
		gin.SetMode(Cfg.Env)
	}
	fmt.Println("now env is set to", Cfg.Env)

	// 解析当前环境的配置文件
	extFile := filepath.Join("config", "config."+env+".yml")
	if extYml, err := ioutil.ReadFile(extFile); err == nil {
		var extCfg Config
		err1 = yaml.Unmarshal(extYml, &extCfg)
		if err1 != nil {
			panic(err1)
		}
		// fmt.Println(extFile, extCfg)
		// 合并配置
		if err = mergo.Merge(&Cfg, extCfg); err != nil {
			fmt.Println("配置文件合并异常", err)
		}
	}

	fmt.Println(Cfg)
}

// 解析并合并对应环境的 yml配置信息
func resloveYml() {

}
