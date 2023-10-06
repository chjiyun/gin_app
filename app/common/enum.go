package common

type ErrType int

const (
	Success ErrType = iota
	Fail
	ParameterError
	IllegalVisit
	UnLogin
	InValidFile
	UnsupportedFormat
	FileNotFound
	UnknownError
	ErrUsernameOrPwd
)

func (t ErrType) String() string {
	switch t {
	case 0:
		return "成功"
	case 1:
		return "失败"
	case 2:
		return "参数错误"
	case 3:
		return "非法访问"
	case 4:
		return "未登录"
	case 5:
		return "无效的文件"
	case 6:
		return "不支持的格式"
	case 7:
		return "文件不存在"
	case 8:
		return "未知错误"
	case 9:
		return "用户名或密码错误"
	default:
		return "未知错误"
	}
}
