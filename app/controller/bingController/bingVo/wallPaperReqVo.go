package bingVo

import (
	"gin_app/app/common"
	"time"
)

type WallPaperReqVo struct {
	common.PageReq
	Keyword        string     `form:"keyword" json:"keyword"`
	ReleaseAtStart *time.Time `form:"releaseAtStart" json:"releaseAtStart"`
	ReleaseAtEnd   *time.Time `form:"releaseAtEnd" json:"releaseAtEnd"`
}
