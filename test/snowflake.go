package test

import (
	"gin_app/util"
	"sync"

	"github.com/gin-gonic/gin"
)

// 并发测试
func Snowflake(c *gin.Context) {
	node := util.SnowFlakeID()
	var wg sync.WaitGroup

	ch := make(chan int64, 10000)
	count := 10000
	wg.Add(count)
	defer close(ch)
	//并发 count个goroutine 进行 snowFlake ID 生成
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := node.Generate()
			ch <- id.Int64()
		}()
	}
	wg.Wait()

	var con = make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := con[id]
		if ok {
			c.JSON(200, gin.H{"repeat id": id})
			return
		}
		// 将 id 作为 key 存入 map
		con[id] = i
	}
	c.JSON(200, con)
}
