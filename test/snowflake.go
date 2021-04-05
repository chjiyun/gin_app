package test

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

// 并发测试
func Snowflake(c *gin.Context) {
	// node := util.SnowFlakeID()

	var wg sync.WaitGroup

	ch := make(chan uint64, 10000)
	count := 10000
	wg.Add(count)
	defer close(ch)
	//并发 count个goroutine 进行 snowFlake ID 生成
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := idgen.NextId()
			ch <- id
		}()
	}
	wg.Wait()

	// var con = make(map[uint64]int)
	// for i := 0; i < count; i++ {
	// 	id := <-ch
	// 	// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
	// 	_, ok := con[id]
	// 	if ok {
	// 		c.JSON(200, gin.H{"repeat id": id})
	// 		return
	// 	}
	// 	// 将 id 作为 key 存入 map
	// 	con[id] = i
	// }

	var arr = make([]uint64, count)
	for i := 0; i < count; i++ {
		id := <-ch
		arr[i] = id
	}

	c.JSON(200, arr)
}
