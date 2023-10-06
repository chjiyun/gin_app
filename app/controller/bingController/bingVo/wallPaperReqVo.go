package bingVo

import (
	"gin_app/app/common"
)

type WallPaperReqVo struct {
	common.PageReq
	Keyword string `form:"keyword" json:"keyword"`
}
