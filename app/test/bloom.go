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
		c.JSON(200, r.Fail("qq missed"))
		return
	}

	rdb := config.RedisDb
	ctx := context.Background()
	// total := 10000 * 10000
	// arr := make([]interface{}, 0, total)
	// var wg sync.WaitGroup

	// for i := 0; i < total; i++ {
	// 	item := util.ToString(i)
	// 	arr = append(arr, item)
	// }
	// fmt.Println("start insert items...")
	// time_start := time.Now()
	// // 并发数
	// asyncCount := 10
	// // 批量插入的数据量
	// size := 1000
	// for i := 0; i < total/(size*asyncCount); i++ {
	// 	wg.Add(asyncCount)
	// 	items := arr[i*asyncCount*size : (i+1)*asyncCount*size]
	// 	for j := 0; j < asyncCount; j++ {
	// 		args := make([]interface{}, 0, size+2)
	// 		args = append(args, "BF.MADD", "qq")
	// 		args = append(args, items[j*size:(j+1)*size]...)
	// 		go func() {
	// 			defer wg.Done()
	// 			_, err := rdb.Do(ctx, args...).Result()
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 		}()
	// 	}
	// 	wg.Wait()
	// 	fmt.Println(i)
	// }
	// duration := time.Since(time_start)
	// fmt.Printf("队列任务耗时：%dms\n", duration/time.Millisecond)

	status, err := rdb.Do(ctx, "BF.EXISTS", "qq", qq).Int()
	if err != nil {
		c.JSON(200, r.Fail(""))
		return
	}

	if status == 1 {
		r.Success("此号码可能存在")
	} else {
		r.Success("此号码不存在")
	}
	c.JSON(200, r)
}
