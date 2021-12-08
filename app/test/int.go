package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Int 数字类型 应用测试
func Int(c *gin.Context) {
	x := "hello"

	now := time.Now()
	fmt.Println(now.UnixNano())

	rand.Seed(now.UnixNano())
	// 生成0到99之间的随机数
	for i := 0; i < 10; i++ {
		n := rand.Intn(100)
		fmt.Printf("Random Number is %d\n", n)
	}

	fmt.Println(RandomString(10, []rune("我是chjiyun")))

	c.JSON(http.StatusOK, x)
}

// 基于随机数生成随机字符串
func RandomString(n int, str ...[]rune) string {
	var letters []rune

	if len(str) == 0 {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	} else {
		letters = str[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
