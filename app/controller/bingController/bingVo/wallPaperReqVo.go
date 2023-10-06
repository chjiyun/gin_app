package bingVo

import (
	"gin_app/app/common"
)

type WallPaperReqVo struct {
	common.PageReq
	Pass    bool   `form:"pass" json:"pass"`
	Keyword string `form:"keyword" json:"keyword"`
}
