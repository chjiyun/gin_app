package common

// Success 成功时的返回体，返回的 Result指针是 r的副本
func (r *Result) Success(msg string, data interface{}) *Result {
	if msg == "" {
		msg = "success"
	}
	r.Code = 200
	r.Msg = msg
	r.Data = data
	return r
}

func (r *Result) SuccessDefault() *Result {
	r.Code = 200
	r.Msg = "success"
	return r
}

func (r *Result) Fail(msg string, data interface{}) *Result {
	if msg == "" {
		msg = "success"
	}
	r.Code = 204
	r.Msg = msg
	r.Data = data
	return r
}

func (r *Result) FailDefault() *Result {
	r.Code = 204
	r.Msg = "failed"
	return r
}

func (r *Result) SetResult(code int, msg string, data interface{}) *Result {
	r.Code = 204
	r.Msg = msg
	r.Data = data
	return r
}
