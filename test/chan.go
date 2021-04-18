package test

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Channel(c *gin.Context) {
	plan := c.Query("p")
	var res interface{}

	switch plan {
	case "1":
		plan1()
	case "2":
		plan2()
	default:
		fmt.Println("非法的参数值")
	}
	c.JSON(200, res)
}

// 方案1
func plan1() (con []interface{}) {
	// 定义一个100个任务
	length := 100
	var tasks = make([]interface{}, length)

	for i := 0; i < length; i++ {
		tasks[i] = i + 1
	}

	ch := make(chan interface{}, 5)
	count := 5
	var wg sync.WaitGroup
	// 把多个组任务分成队列，组内 count个任务进行并发
	for i := 0; i < length/count; i++ {
		arr := tasks[i*count : (i+1)*count]
		wg.Add(count)
		for j := 0; j < count; j++ {
			go func(j0 int) {
				defer wg.Done()
				// 模拟函数执行时间
				time.Sleep(10 * time.Millisecond)
				ch <- arr[j0]
				// fmt.Println("push >>", arr[j0])
			}(j)
		}
		wg.Wait()

		for j := 0; j < count; j++ {
			val := <-ch
			// fmt.Println("caught <<", val)
			con = append(con, val)
		}
	}
	return
}

// 方案2
func plan2() (con []interface{}) {
	// waitGroup 对象不是一个引用类型，在通过函数传值时需使用地址
	var wg sync.WaitGroup
	wg.Add(2)
	// make返回的是引用类型本身,而new返回的是指向类型的指针
	// 无缓冲通道，推入一个值即开始阻塞，待取出后继续推入
	ch := make(chan interface{})

	go thrower(ch, &wg)
	go catcher(ch, &wg)

	wg.Wait()
	return
}

func thrower(c chan interface{}, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		var arr []int
		for j := i*5 + 1; j <= 5*(i+1); j++ {
			arr = append(arr, j)
		}
		c <- arr
		fmt.Println("push >>", arr)
	}
	wg.Done()
}

func catcher(c chan interface{}, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		val := <-c
		// 模拟sql更新等耗时操作
		time.Sleep(100 * time.Millisecond)
		fmt.Println("caught <<", val)
	}
	wg.Done()
}
