package test

import (
	"context"
	"gin_app/app/result"
	"gin_app/config"

	"github.com/gin-gonic/gin"
)

// BloomFilter 布隆过滤器去重实例
func BloomFilter(c *gin.Context) {
	r := result.New()
	qq := c.Query("qq")
	if qq == "" {
		r.Fail("qq missed", nil)
		c.JSON(200, r)
		return
	}

	redis := config.RedisDb
	ctx := context.Background()

	status, err := redis.Do(ctx, "BF.EXISTS", "qq", qq).Result()
	if err != nil {
		r.Fail("", err)
		return
	}
	r.SetData(status)
	c.JSON(200, r)
}
