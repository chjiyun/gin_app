package fileVo

import "gin_app/app/common"

type ThumbRespVo struct {
	common.BaseModel
	Uid    uint64 `gorm:"not null;comment:雪花id" json:"uid"`
	FileId string `gorm:"not null;comment:file_id" json:"file_id"`
	Ext    string `gorm:"size:20;comment:文件后缀" json:"ext"`
	Width  uint   `gorm:"comment:宽度" json:"width"`
	Height uint   `gorm:"comment:高度" json:"height"`
	Name   string `gorm:"size:120;comment:原始文件名" json:"name"`
	Size   uint   `gorm:"comment:文件大小（Byte）" json:"size"`
}
