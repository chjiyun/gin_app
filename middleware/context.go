package middleware

import (
	"context"
	"gin_app/config"
	"time"

	"github.com/gin-gonic/gin"
)

// 初始化后的 *gorm.DB 放到 gin.context
func SetDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置超时 Context
		timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.Set("DB", config.DB.WithContext(timeoutContext))
		c.Next()
	}
}
