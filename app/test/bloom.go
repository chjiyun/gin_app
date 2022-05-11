package test

import (
	"context"
	"fmt"
	"gin_app/app/result"
	"gin_app/app/util"
	"gin_app/config"
	"sync"
	"time"

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

	arr := make([]string, 0, 1e8)
	var wg sync.WaitGroup
	// 插入过滤器的错误集合
	results := make([]string, 0, 1e4)

	for i := 0; i < 1e8; i++ {
		item := util.ToString(i)
		arr = append(arr, item)
	}
	time_start := time.Now()
	// 10000个任务排队
	for i := 0; i < 1e4; i++ {
		items := arr[i*1e4 : (i+1)*1e4]
		wg.Add(1e4)
		// 10000个协程=并发量
		for j := 0; j < len(items); j++ {
			go func(s string) {
				defer wg.Done()
				_, err := redis.Do(ctx, "BF.ADD", "qq", s).Result()
				if err != nil {
					results = append(results, s)
				}
			}(items[j])
		}
		wg.Wait()
	}
	duration := time.Since(time_start)
	fmt.Printf("队列任务耗时：%d", duration)

	r.SetData(results)
	c.JSON(200, r)
}
