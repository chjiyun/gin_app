package bingVo

import (
	"mime/multipart"
	"time"
)

// 接收post form参数
type WallPaperCreateReqVo struct {
	File      *multipart.FileHeader `form:"file" json:"file"`
	Desc      string                `form:"desc" json:"desc"`
	ReleaseAt *time.Time            `form:"releaseAt" json:"releaseAt"`
}
