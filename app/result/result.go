package result

import (
	"gin_app/app/common"
	"reflect"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// New return a new Result instance
func New() Result {
	r := Result{}
	r.Code = 0
	r.Msg = common.Success.String()
	return r
}

// Success 成功时的返回体，返回的 Result指针是 r的副本
func (r Result) Success(data interface{}) Result {
	r.Code = 0
	r.Data = data
	return r
}

// Fail 失败时的输出
func (r Result) Fail(msg string) Result {
	if msg == "" {
		msg = common.Fail.String()
	}
	r.Code = 1
	r.Msg = msg
	return r
}

// FailType 指定错误类型（枚举）
func (r Result) FailType(errType common.ErrType) Result {
	r.Code = 1
	r.Msg = errType.String()
	return r
}

// FailErr 直接输出error
func (r Result) FailErr(err error) Result {
	r.Code = 1
	r.Msg = err.Error()
	return r
}

func (r Result) SetData(data interface{}) {
	r.Data = data
}

func (r Result) SetCode(code int) Result {
	r.Code = code
	return r
}

// handleData 处理data特殊类型
func handleData(data, d interface{}) {
	// 取出指针的值
	dataVal := reflect.ValueOf(data).Elem().Interface()
	dVal := reflect.ValueOf(d).Elem()
	switch dataVal.(type) {
	case error:
		dVal.Set(reflect.ValueOf(dataVal.(error).Error()))
	default:
		return
	}
}
