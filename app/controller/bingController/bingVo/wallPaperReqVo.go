package bingVo

import (
	"gin_app/app/common"
)

type WallPaperReqVo struct {
	common.PageReq
	Status  string `form:"status" json:"status"`
	Keyword string `form:"keyword" json:"keyword"`
}
