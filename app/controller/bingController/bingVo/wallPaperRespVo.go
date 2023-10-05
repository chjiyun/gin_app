package bingVo

import (
	"gin_app/app/common"
	"time"
)

// WallPaperRespVo 壁纸相应Vo
type WallPaperRespVo struct {
	common.BaseModel
	FileId    string    `gorm:"not null;comment:外键" json:"file_id"`
	Desc      string    `gorm:"size:500;comment:备注描述" json:"desc"`
	ReleaseAt time.Time `gorm:"type:date;comment:发布日期" json:"releaseAt"`
	Width     uint      `gorm:"comment:宽度" json:"width"`
	Height    uint      `gorm:"comment:高度" json:"height"`
	Ext       string    `gorm:"size:20;comment:原始文件类型" json:"type"`
	Name      string    `gorm:"size:120;comment:原始文件名" json:"name"`
	ThumbId   string    `json:"thumb_id"`
	CreatedAt time.Time `json:"-"`
}
