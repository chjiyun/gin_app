package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
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

// ParseInt 字符串数字转整型
func ParseInt(s string, base int, bitSize int) {
	// v, err := strconv.Atoi(s)

}

// ToJson 转成json字符串
func ToJson(m interface{}) string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// SendFormData 以multipart/form-data格式发送文件,
// fileField: file字段名，data: form键值对
func SendFormData(url string, fileField string, data map[string]interface{}) (*http.Response, error) {
	filename, ok := data["filename"].(string)
	if !ok || len(filename) == 0 {
		return nil, errors.New("filename is not exist")
	}
	fileByte, ok := data[fileField].(*[]byte)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s []byte pointer is not exist", fileField))
	}
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)
	for key, value := range data {
		if key == fileField || key == "filename" {
			continue
		}
		bodyWriter.WriteField(key, value.(string))
	}
	fileWriter, err := bodyWriter.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, err
	}
	_, err = fileWriter.Write(*fileByte)
	// os.File 是 io.reader的实现
	// _, err = io.Copy(fileWriter, data[fileField].(*os.File))
	if err != nil {
		return nil, err
	}
	// 完成所有内容设置后，一定要关闭 Writer，否则，请求体会缺少结束边界
	err = bodyWriter.Close()
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, bodyWriter.FormDataContentType(), buf)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// handleData data传指针，处理error错误
func HandleData(data interface{}) interface{} {
	// 获取指针指向的值
	ptrVal := reflect.Indirect(reflect.ValueOf(data)).Interface()
	switch ptrVal.(type) {
	case error:
		return ptrVal.(error).Error()
	default:
		// 原路返回
		return ptrVal
	}
}

// TypeOf 获取变量类型
func TypeOf(data interface{}) string {
	if data == nil {
		return "nil"
	}
	return reflect.TypeOf(data).Kind().String()
}

// SnakeString 驼峰转蛇形名
// XxYy to xx_yy , XxYY to xx_y_y
func SnakeString(s string) string {
	length := len(s)
	data := make([]byte, 0, len(s)*2)
	for i := 0; i < length; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' {
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
