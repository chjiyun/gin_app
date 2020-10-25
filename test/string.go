package test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// String 数组和切片
func String(c *gin.Context) {
	// 一个 UTF8 编码的字符可能会占多个字节，比如汉字就需要 3~4 个字节来存储
	// 	golang的string是以UTF-8编码的,而UTF-8是一种1-4字节的可变长字符集，每个字符可用1-4字节 来表示
	// 使用下标方式s[i]访问字符串s，s[i]是UTF-8编码后的一个字节(uint8)，即按字节遍历
	// 使用for i,v := range s 方式访问s，i是字符串下标编号，v是对应的字符值(int32=rune)，即按字符遍历
	x := "hello"
	// 将 string 转为 rune slice（代表一个 UTF-8 字符，rune 类型等价于 int32 类型）
	xRunes := []rune(x)
	xRunes[0] = '我'
	fmt.Println(xRunes)
	x = string(xRunes)

	// 截取字符串
	tracer := "死神来了，死神bye bye"
	strSlice := []rune(tracer)
	fmt.Println(string(strSlice[5:]))
	fmt.Println(strings.Trim(tracer, "bye"))
	fmt.Println(strings.TrimSuffix(tracer, "bye"))

	// 查找子串的位置(字节)
	fmt.Println(strings.IndexRune(tracer, rune('了')))

	fmt.Println(UnicodeIndex(tracer, "来了"))

	// 插入字符串
	strSlice = append(strSlice[:5], append([]rune("where? "), strSlice[5:]...)...)
	fmt.Println(string(strSlice))

	// 拼接字符串
	var builder bytes.Buffer
	builder.WriteString("Hello,")
	builder.WriteString("world!")
	fmt.Println(builder.String())

	c.JSON(http.StatusOK, x)
}

// 查找子串位置（偏移量）
func UnicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}

	return result
}
