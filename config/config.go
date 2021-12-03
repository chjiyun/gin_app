package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	Basedir    string
	Name       string       `yaml:"name"`
	Env        string       `yaml:"env"`
	Server     Server       `yaml:"server"`
	Redis      Redis        `yaml:"redis"`
	Datasource []Datasource `yaml:"datasource"`
	Log        Log          `yaml:"log"`
}
type Server struct {
	Port string `yaml:"port"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type Datasource struct {
	Model   string `yaml:"model"`
	Dialect string `yaml:"dialect"`
	Dsn     string `yaml:"dsn"`
}
type Log struct {
	Filename string `yaml:"filename"`
	Filepath string `yaml:"filepath"`
}

// 配置信息缓存
var Cfg Config

// redis 实例
var RedisDb *redis.Client

// db实例指针
var DB *gorm.DB

// 初始化 config 配置
func Init() {
	// 设置根目录
	Cfg.Basedir, _ = filepath.Abs(".")

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
	switch env {
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

	fmt.Println("merge Config:", Cfg)

	redisInit()

	dbInit()
}

// 解析并合并对应环境的 yml配置信息
func resloveYml() {

}

// 初始化 redis
func redisInit() {
	if Cfg.Redis.Addr == "" {
		return
	}
	// 定义一个 reids的 options 结构体
	var options redis.Options
	// 拷贝结构体
	err := copier.CopyWithOption(&options, &Cfg.Redis, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		panic(err)
	}
	RedisDb := redis.NewClient(&options)
	_, err = RedisDb.Ping(RedisDb.Context()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("redis is connected to %s\n", options.Addr)
}

// 初始化 db
func dbInit() {
	length := len(Cfg.Datasource)
	if length == 0 {
		return
	}
	var err error
	// 设置sql日志级别
	logMode := logger.Warn
	if Cfg.Env == gin.DebugMode {
		logMode = logger.Info
	}
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: Cfg.Datasource[0].Dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("数据源：%s 已经连接\n", Cfg.Datasource[0].Model)
	if length == 1 {
		return
	}
	// init other db
	var dr *dbresolver.DBResolver
	for i, ds := range Cfg.Datasource[1:] {
		dr_cfg := dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open(ds.Dsn)},
		}
		if i == 0 {
			dr = dbresolver.Register(dr_cfg)
		} else {
			dr.Register(dr_cfg)
		}
		fmt.Printf("数据源：%s 已经连接\n", ds.Model)
	}
	DB.Use(dr)
}
