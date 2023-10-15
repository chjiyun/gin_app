package bingVo

import (
	"time"
)

// 接收post form参数
type WallPaperUpdateReqVo struct {
	ID        string     `form:"id" json:"id" binding:"required"`
	FileId    string     `form:"file_id" json:"file_id"`
	Desc      string     `form:"desc" json:"desc"`
	ReleaseAt *time.Time `form:"releaseAt" json:"releaseAt"`
}
