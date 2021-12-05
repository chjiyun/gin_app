package middleware

import (
	"gin_app/config"
	"gin_app/util"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	logFilePath := config.Cfg.Log.Filepath
	logFileName := config.Cfg.Log.Filename
	// 生成win下的日志文件夹，相对路径
	if runtime.GOOS == "windows" && filepath.VolumeName(logFilePath) == "" {
		logFilePath = filepath.Join(config.Cfg.Basedir, logFilePath)
		// 文件夹不存在则创建
		if !util.CheckFileIsExist(logFilePath) {
			err := os.Mkdir(logFilePath, 0666)
			if err != nil {
				panic(err)
			}
		}
		p := filepath.Join(logFilePath, logFileName+".log")
		if !util.CheckFileIsExist(p) {
			os.Create(p)
		}
	}

	// 日志文件
	fileName := filepath.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName+".log", os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	// 实例化
	// logger := logrus.New()
	// 设置输出
	if config.Cfg.Env == gin.DebugMode {
		w := io.MultiWriter(file, os.Stdout)
		config.Log.SetOutput(w)
		gin.DefaultWriter = w
	} else {
		config.Log.SetOutput(file)
	}
	// 设置日志级别
	config.Log.SetLevel(logrus.DebugLevel)
	// 输出行号
	// config.Log.SetReportCaller(true)
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 新增 Hook
	config.Log.AddHook(lfHook)
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		config.Log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

// // 日志记录到 MongoDB
// func LoggerToMongo() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	}
// }

// // 日志记录到 ES
// func LoggerToES() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	}
// }

// // 日志记录到 MQ
// func LoggerToMQ() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	}
// }
