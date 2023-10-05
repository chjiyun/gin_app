package bingVo

import (
	"gin_app/app/common"
	"time"
)

type BingPageReqVo struct {
	common.PageReq
	StartTime time.Time `form:"start_time" json:"start_time"`
	EndTime   time.Time `form:"end_time" json:"end_time"`
}
