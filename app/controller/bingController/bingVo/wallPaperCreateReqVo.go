package bingVo

import (
	"time"
)

// 接收post form参数
type WallPaperCreateReqVo struct {
	FileId    string     `form:"file_id" json:"file_id" binding:"required"`
	Desc      string     `form:"desc" json:"desc"`
	ReleaseAt *time.Time `form:"releaseAt" json:"releaseAt"`
}
