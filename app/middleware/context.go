package middleware

import (
	"context"
	"gin_app/config"
	"time"

	"github.com/gin-gonic/gin"
)

// context传递
func SetContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置超时 Context
		timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		c.Set("DB", config.DB.WithContext(timeoutContext))
		// Add a context to the log entry.
		c.Set("Logger", config.Logger.WithContext(context.Background()))
		c.Next()
	}
}
