package bingVo

import (
	"time"
)

type WallPaperAuditReqVo struct {
	ID        uint64     `json:"id" binding:"required"`
	Status    string     `json:"status" binding:"oneof=1 2"`
	Remarks   string     `json:"remarks"`
	ReleaseAt *time.Time `json:"releaseAt"`
}
