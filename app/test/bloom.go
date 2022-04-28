package test

import (
	"fmt"
	"gin_app/config"
	"log"

	"github.com/gin-gonic/gin"
)

// BloomFilter 布隆过滤器去重实例
func BloomFilter(c *gin.Context) {
	redis := config.RedisDb
	cmd := redis.Do(c, "BF.EXISTS", "qq", "123456")
	status, err := cmd.Int()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(status)
}
