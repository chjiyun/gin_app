package model

import (
	"gin_app/app/common"
	"time"
)

type Bing struct {
	common.BaseModel
	FileId    string    `gorm:"not null;comment:外键" json:"file_id"`
	Url       string    `gorm:"size:255;comment:原bing图片链接" json:"url"`
	Hsh       string    `gorm:"size:32;comment:图片唯一hash_id" json:"hsh"`
	Desc      string    `gorm:"size:500;comment:描述" json:"desc"`
	ReleaseAt time.Time `gorm:"type:date;comment:发布日期" json:"releaseAt"`
	IsBing    bool      `gorm:"not null;comment:是否来自bing壁纸" json:"isBing"`
	Status    string    `gorm:"not null;comment:审核通过状态" json:"status"`
	Remarks   string    `gorm:"size:255;comment:备注" json:"remarks"`
	File      File      `json:"file"`
}
