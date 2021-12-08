package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Arr 数组和切片
func Arr(c *gin.Context) {
	// x := "hello"
	// 初始化数组
	q := [...]int{1, 2, 3}

	a := q[0:]                 // 从数组生成切片
	a = append(a, 4)           // 追加一个元素
	a = append([]int{0}, a...) // 在一个切片追加多个元素，相当于在开头插入一个元素
	a = append([]int{-2, -1}, a...)
	a = append(a[:2], append([]int{0}, a[3:]...)...) // 在第2个位置插入元素
	fmt.Println(a)

	obj := []interface{}{"s", a, "e"}
	b, _ := json.Marshal(obj)
	fmt.Println(string(b))

	c.JSON(http.StatusOK, obj)
}
