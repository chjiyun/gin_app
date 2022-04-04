// 接口response 定义与配置
package result

import (
	"reflect"
)

type Result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}
type message struct {
	code int
	msg  string
}

var ResultMap = map[string]message{
	"success": {
		code: 200,
		msg:  "成功",
	},
	"parameterError": {
		code: 202,
		msg:  "参数错误",
	},
	"illegalVisit": {
		code: 203,
		msg:  "非法访问",
	},
	"fail": {
		code: 204,
		msg:  "失败",
	},
	"unLogin": {
		code: 205,
		msg:  "未登录",
	},
	"inValidFile": {
		code: 206,
		msg:  "无效的文件",
	},
	"notFound": {
		code: 207,
		msg:  "file not found",
	},
	"serverError": {
		code: 208,
		msg:  "服务内部错误",
	},
}

// New return a new Result instance
func New() *Result {
	r := &Result{}
	res := ResultMap["success"]
	r.Code = res.code
	r.Msg = res.msg
	return r
}

// Success 成功时的返回体，返回的 Result指针是 r的副本
func (r *Result) Success(msg string, data interface{}) *Result {
	res := ResultMap["success"]
	if msg == "" {
		msg = "success"
	}
	r.Code = res.code
	r.Msg = msg
	r.Data = data
	return r
}

func (r *Result) SuccessDefault() *Result {
	res := ResultMap["success"]
	r.Code = res.code
	r.Msg = res.msg
	return r
}

// Fail 失败时的输出
func (r *Result) Fail(msg string, err error) *Result {
	res := ResultMap["fail"]
	if msg == "" {
		msg = res.msg
	}
	r.Code = res.code
	r.Msg = msg
	r.Error = err
	return r
}

func (r *Result) FailDefault() *Result {
	res := ResultMap["fail"]
	r.Code = res.code
	r.Msg = res.msg
	return r
}

func (r *Result) SetData(data interface{}) {
	r.Data = data
}

// SetResult 自定义code msg
func (r *Result) SetResult(res message, msg string) *Result {
	r.Code = res.code
	r.Msg = res.msg
	if msg != "" {
		r.Msg = msg
	}
	return r
}

// SetError 错误输出
func (r *Result) SetError(err error) {
	r.Error = err
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
