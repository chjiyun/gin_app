// 接口response 定义与配置
package common

type message struct {
	code int
	msg  string
}

// const (
// 	Success = iota+1
// 	ParameterError
// 	IllegalVisit
// 	Fail
// 	UnLogin
// 	InValidFile
// 	NotFound
// )
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

func (r *Result) Fail(msg string, data interface{}) *Result {
	res := ResultMap["fail"]
	if msg == "" {
		msg = res.msg
	}
	r.Code = res.code
	r.Msg = msg
	r.Data = data
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
