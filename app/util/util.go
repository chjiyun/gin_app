package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
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
func Basename(filename string) string {
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '.' {
			return filename[:i]
		}
	}
	return filename
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

// ToString 各种类型转string
// 整数转换为10进制的字符串
func ToString(v interface{}) string {
	t := reflect.TypeOf(v)
	var s string
	switch t.Kind() {
	case reflect.Int:
		s = strconv.FormatInt(int64(v.(int)), 10)
	case reflect.Int64:
		s = strconv.FormatInt(int64(v.(int64)), 10)
	case reflect.Int16:
		s = strconv.FormatInt(int64(v.(int16)), 10)
	case reflect.Int8:
		s = strconv.FormatInt(int64(v.(int8)), 10)
	case reflect.Uint:
		s = strconv.FormatUint(uint64(v.(uint)), 10)
	case reflect.Uint64:
		s = strconv.FormatUint(v.(uint64), 10)
	case reflect.Uint16:
		s = strconv.FormatUint(uint64(v.(uint16)), 10)
	case reflect.Uint8:
		s = strconv.FormatUint(uint64(v.(uint8)), 10)
	case reflect.Bool:
		s = strconv.FormatBool(v.(bool))
	case reflect.Float32:
		// 默认以(-ddd.dddd, no exponent)格式转化浮点数
		s = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case reflect.Float64:
		s = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case reflect.Map, reflect.Struct, reflect.Slice:
		s = ToJson(v)
	default:
		fmt.Printf("type %s is not support, use fmt.Sprintf instead", t.Kind())
	}
	return s
}

func ParseInt() {
}

// ToJson 转成json字符串
func ToJson(m interface{}) string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}
