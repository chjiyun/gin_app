package userVo

import "gin_app/app/common"

// UserPageReqVo 用户分页列表请求Vo
type UserPageReqVo struct {
	common.PageReq
	Keyword string `form:"keyword" json:"keyword"`
}
