package middleware

import (
	"gin_app/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
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
		config.Logger.WithFields(logrus.Fields{
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
