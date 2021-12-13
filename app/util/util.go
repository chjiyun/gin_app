package util

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Call 反射动态调用函数
func Call(m map[string]interface{}, fnName string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[fnName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return nil, err
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return result, nil
}

// CheckFileIsExist 检查文件或目录是否存在
func CheckFileIsExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// SnowFlakeID 产生一个 snowflake node
// func SnowFlakeID() *snowflake.Node {
// 	// Create a new Node with a Node number of 1
// 	node, err := snowflake.NewNode(1)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return node
// }

// TimeCost 函数耗时统计
func TimeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}

// SetDefault 给变量设置默认值
func SetDefault(v, _default interface{}) {
	v1 := reflect.ValueOf(v).Elem()
	v2 := reflect.ValueOf(_default)
	// 初始化完成的map 和 数组 不会被覆盖
	if v1.IsZero() {
		v1.Set(v2)
	}
}

// UpperFirst 字符串首字母大写
func UpperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LowerFirst 字符串首字母小写
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Basename 获取文件基础名，禁止含非1字节字符
func Basename(str string) string {
	for i := len(str) - 1; i > 0; i-- {
		if str[i] == '.' {
			return str[:i]
		}
	}
	return ""
}

// WriteString 拼接字符串
func WriteString(str ...string) string {
	if len(str) == 1 {
		return str[0]
	}
	var b strings.Builder
	for _, s := range str {
		b.WriteString(s)
	}
	return b.String()
}

// 读取目录下的特定后缀文件基础名，首字母大写（ex: app.go -> App）；
// fileExt为空 返回文件名
func GetFileBasename(dirname string, fileExt []string) []string {
	var names []string
	fileInfo, _ := os.ReadDir(dirname)
	if len(fileInfo) == 0 {
		return names
	}
	var str string
	switch len(fileExt) {
	case 0:
		for _, f := range fileInfo {
			names = append(names, f.Name())
		}
		return names
	case 1:
		str = fileExt[0]
	default:
		str = strings.Join(fileExt, "|")
	}
	str = WriteString(`^\w+\.(`, str, ")$")
	re := regexp.MustCompile(str)
	for _, file := range fileInfo {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		// 匹配 .go结尾的文件
		if re.MatchString(name) {
			basename := Basename(name)
			basename = UpperFirst(basename)
			names = append(names, basename)
		}
	}
	return names
}
