package bingVo

import (
	"time"
)

// 接收post form参数
type WallPaperCreateReqVo struct {
	FileId    string     `json:"fileId" binding:"required"`
	Desc      string     `form:"desc" json:"desc"`
	ReleaseAt *time.Time `form:"releaseAt" json:"releaseAt"`
}
