package config

import (
	"fmt"
	"gin_app/app/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
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
	Jwt        JwtConfig    `yaml:"jwt"`
}
type Server struct {
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
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
	Filename      string `yaml:"filename"`
	Filepath      string `yaml:"filepath"`
	LongQueryTime int    `yaml:"long_query_time"`
	Level         string `yaml:"level"`
}

type JwtConfig struct {
	Issuer    string `json:"issuer"`
	Audience  string `json:"audience"`
	Expires   int    `json:"expires"`
	SecretKey string `json:"secret_key"`
	Whitelist string `json:"whitelist"`
	Refresh   int    `json:"refresh"`
}

type SqlWriter struct {
	log *zap.SugaredLogger
}

var Cfg Config

var RedisDb *redis.Client

var DB *gorm.DB

//var Logger = logrus.New()

var SugarLog *zap.SugaredLogger

// 初始化 config 配置
func Init() {
	// 解析默认基础配置文件
	filename := filepath.Join("config", "config.yml")
	yml, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err1 := yaml.Unmarshal(yml, &Cfg)
	if err1 != nil {
		panic(err1)
	}
	// fmt.Println(filename, Cfg)
	// 设置一些默认值
	Cfg.Basedir, _ = filepath.Abs(".")
	if Cfg.Name == "" {
		Cfg.Name = filepath.Base(Cfg.Basedir)
	}

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
	if extYml, err := os.ReadFile(extFile); err == nil {
		var extCfg Config
		err1 = yaml.Unmarshal(extYml, &extCfg)
		if err1 != nil {
			panic(err1)
		}
		// fmt.Println(extFile, extCfg)
		// 合并配置
		err1 = copier.CopyWithOption(&Cfg, &extCfg, copier.Option{DeepCopy: true, IgnoreEmpty: true})
		if err1 != nil {
			panic(err1)
		}
	}
	if err = os.MkdirAll(filepath.Join(Cfg.Basedir, "files/temp"), 0666); err != nil {
		panic(err)
	}

	// fmt.Println("merge Config:", Cfg)
	// 按顺序执行
	//LogInit()
	zapLogInit()
	redisInit()
	dbInit()
}

// 解析并合并对应环境的 yml配置信息
// func resloveYml() {

// }

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
	// options.PoolSize = 1000
	// options.MinIdleConns = 100
	RedisDb = redis.NewClient(&options)
	_, err = RedisDb.Ping(RedisDb.Context()).Result()
	if err != nil {
		panic(err)
	}
	SugarLog.Infof("redis is connected to %s", options.Addr)
}

// 初始化 db
func dbInit() {
	length := len(Cfg.Datasource)
	if length == 0 {
		return
	}
	var err error
	// 设置sql日志打印级别  开发环境下才会输出sql到日志
	logMode := logger.Warn
	if Cfg.Env == gin.DebugMode {
		logMode = logger.Info
		for i, ds := range Cfg.Datasource {
			if ds.Dsn == "" {
				continue
			}
			key, err := os.ReadFile("hashkey.txt")
			if err != nil {
				panic(err)
			}
			Cfg.Datasource[i].Dsn = util.Decrypt(ds.Dsn, key)
		}
	}
	// sql记录到日志
	sqlLogger := logger.New(&SqlWriter{log: SugarLog}, logger.Config{
		SlowThreshold:             time.Duration(Cfg.Log.LongQueryTime) * time.Millisecond,
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: true,
	})

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  Cfg.Datasource[0].Dsn,
		PreferSimpleProtocol: true, //禁用 prepared statement 缓存
		//Conn:                 conn,
	}), &gorm.Config{
		Logger: sqlLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	// 若服务端设置保活了机制 则必须小于 tcp_keepalives_idle=300
	db.SetConnMaxLifetime(time.Minute * 5)

	SugarLog.Infof("datasource: %s has been connected", Cfg.Datasource[0].Model)

	if length == 1 {
		return
	}
	// init other db
	var dr *dbresolver.DBResolver
	for i, ds := range Cfg.Datasource[1:] {
		drCfg := dbresolver.Config{
			Sources: []gorm.Dialector{postgres.Open(ds.Dsn)},
		}
		if i == 0 {
			dr = dbresolver.Register(drCfg)
		} else {
			dr.Register(drCfg)
		}
		SugarLog.Infof("datasource: %s has been connected", ds.Model)
	}
	DB.Use(dr)
}

// logrus实例初始化
//func LogInit() {
//	logFilePath := Cfg.Log.Filepath
//	logFileName := Cfg.Log.Filename
//	if logFileName == "" {
//		logFileName = Cfg.Name
//	}
//
//	if logFilePath == "" {
//		switch runtime.GOOS {
//		case "darwin", "windows":
//			// 没配置path就在项目根目录创建文件夹
//			logFilePath = filepath.Join(Cfg.Basedir, "logs")
//		default:
//			// log目录下再加同名项目文件夹
//			logFilePath = filepath.Join("/root/logs", Cfg.Name)
//		}
//	}
//	// 创建路径中缺失的文件夹
//	if !util.CheckFileIsExist(logFilePath) {
//		err := os.MkdirAll(logFilePath, 0666)
//		if err != nil {
//			panic(err)
//		}
//	}
//	baseLog := filepath.Join(logFilePath, logFileName+".log")
//
//	fileName := filepath.Join(logFilePath, logFileName)
//	// 写入文件（0660：其他用户的权限）
//	file, err := os.OpenFile(baseLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
//	if err != nil {
//		panic(err)
//	}
//	// 设置输出
//	if Cfg.Env == gin.DebugMode {
//		w := io.MultiWriter(os.Stdout, file)
//		Logger.Out = w
//		gin.DefaultWriter = w
//		gin.DefaultErrorWriter = io.MultiWriter(file, os.Stderr)
//	} else {
//		gin.DisableConsoleColor()
//		Logger.SetOutput(file)
//	}
//	// 设置日志级别
//	Logger.SetLevel(logrus.DebugLevel)
//	// 输出行号
//	// Logger.SetReportCaller(true)
//	// 设置 rotatelogs
//	logWriter, err := rotatelogs.New(
//		// 分割后的文件名称
//		fileName+".%Y%m%d.log",
//		// 生成软链，指向最新日志文件
//		rotatelogs.WithLinkName(fileName),
//		// 设置最大保存时间(7天)
//		rotatelogs.WithMaxAge(30*24*time.Hour),
//		// 设置日志切割时间间隔(1天)
//		rotatelogs.WithRotationTime(24*time.Hour),
//	)
//	if err != nil {
//		panic(err)
//	}
//	writeMap := lfshook.WriterMap{
//		logrus.InfoLevel:  logWriter,
//		logrus.FatalLevel: logWriter,
//		logrus.DebugLevel: logWriter,
//		logrus.WarnLevel:  logWriter,
//		logrus.ErrorLevel: logWriter,
//		logrus.PanicLevel: logWriter,
//	}
//	lfHook := lfshook.NewHook(writeMap, &nested.Formatter{
//		TimestampFormat: "2006-01-02 15:04:05",
//		// PrettyPrint:     true,
//		HideKeys: true,
//		NoColors: true,
//	})
//	// 新增 Hook
//	Logger.AddHook(lfHook)
//}

// 日志自定义配置
func zapLogInit() {
	logFilePath := Cfg.Log.Filepath
	logFileName := Cfg.Log.Filename
	if logFileName == "" {
		logFileName = Cfg.Name
	}

	if logFilePath == "" {
		switch runtime.GOOS {
		case "darwin", "windows":
			// 没配置path就在项目根目录创建文件夹
			logFilePath = filepath.Join(Cfg.Basedir, "logs")
		default:
			// log目录下再加同名项目文件夹
			logFilePath = filepath.Join("/root/logs", Cfg.Name)
		}
	}
	// 创建路径中缺失的文件夹
	if !util.CheckFileIsExist(logFilePath) {
		err := os.MkdirAll(logFilePath, 0666)
		if err != nil {
			panic(err)
		}
	}
	filename := filepath.Join(logFilePath, logFileName+".log")

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	levelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}
	callerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		path := util.WriteString("[", caller.TrimmedPath(), "]")
		enc.AppendString(path)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = levelEncoder
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeCaller = callerEncoder
	encoderConfig.ConsoleSeparator = " "
	fmt.Println(filename)
	logWriter := &lumberjack.Logger{
		Filename:  filename,
		MaxSize:   20,
		MaxAge:    30,
		LocalTime: true,
		Compress:  false,
	}
	var syncWriter zapcore.WriteSyncer
	if Cfg.Env == gin.DebugMode {
		syncWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logWriter))
	} else {
		syncWriter = zapcore.AddSync(logWriter)
		gin.DisableConsoleColor()
	}
	gin.DefaultWriter = syncWriter

	var level zapcore.Level
	switch Cfg.Log.Level {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "Error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.WarnLevel
	}
	zapCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), syncWriter, level)
	SugarLog = zap.New(zapCore, zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}

// Printf 实现gorm/logger.Writer接口
func (m *SqlWriter) Printf(format string, v ...interface{}) {
	m.log.Infof(format, v...)
}
