package adminVo

import (
	"gin_app/app/common"
	"time"
)

type WallPaperReqVo struct {
	common.PageReq
	Status         string     `form:"status" json:"status"`
	Keyword        string     `form:"keyword" json:"keyword"`
	ReleaseAtStart *time.Time `form:"releaseAtStart" json:"releaseAtStart"`
	ReleaseAtEnd   *time.Time `form:"releaseAtEnd" json:"releaseAtEnd"`
}

type WallPaperRespVo struct {
	common.BaseModel
	FileId    string    `gorm:"not null;comment:外键" json:"file_id"`
	Desc      string    `gorm:"size:500;comment:描述" json:"desc"`
	Remarks   string    `gorm:"size:255;comment:备注" json:"remarks"`
	Status    string    `gorm:"not null;comment:审核通过状态" json:"status"`
	ReleaseAt time.Time `gorm:"type:date;comment:发布日期" json:"releaseAt"`
	Width     uint      `gorm:"comment:宽度" json:"width"`
	Height    uint      `gorm:"comment:高度" json:"height"`
	Ext       string    `gorm:"size:20;comment:原始文件类型" json:"type"`
	Name      string    `gorm:"size:120;comment:原始文件名" json:"name"`
	ThumbId   string    `json:"thumb_id"`
}

type WallPaperAuditReqVo struct {
	ID        uint64     `json:"id" binding:"required"`
	Status    string     `json:"status" binding:"oneof=1 2"`
	Remarks   string     `json:"remarks"`
	ReleaseAt *time.Time `json:"releaseAt"`
}
