package test

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Cat struct {
	Color string
	Name  string
	Age   int
}
type BlackCat struct {
	cat Cat
	Cat // 嵌入Cat, 类似于派生
	C   *Cat
}

// “构造基类”
func NewCat(name string) *Cat {
	return &Cat{
		Name: name,
		Age:  20,
	}
}

// “构造子类”
func NewBlackCat(color string) *BlackCat {
	cat := &BlackCat{} // 同时会实例化Cat，自由访问 Cat 的所有成员
	cat.Color = color
	cat.cat.Color = color
	// 指针结构体初始化
	// cat.C = &Cat{}
	// cat.C.Color = color
	return cat
}

// 模拟面向对象的继承
func StructNest(c *gin.Context) {
	fmt.Println(NewCat("Tom"))        // &{ Tom 20}
	fmt.Println(NewBlackCat("Green")) // &{{green  0} {green  0} <nil>}
}
