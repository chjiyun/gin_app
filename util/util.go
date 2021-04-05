package util

import (
	"errors"
	"os"
	"reflect"
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

// CheckFileIsExist 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
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
