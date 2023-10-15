package bingVo

import (
	"time"
)

// 接收post form参数
type WallPaperUpdateReqVo struct {
	ID        string     `json:"id" binding:"required"`
	FileId    string     `json:"fileId"`
	Desc      string     `form:"desc" json:"desc"`
	ReleaseAt *time.Time `form:"releaseAt" json:"releaseAt"`
}
