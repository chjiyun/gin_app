package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// soft_delete 此插件还支持混合模式，删除同时更新删除时间
type BaseModel struct {
	ID        uint                  `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag;not null;default:0" json:"-"`
}

type Bing struct {
	BaseModel
	FileId    uint      `gorm:"not null;comment:外键" json:"file_id"`
	Url       string    `gorm:"size:255;comment:原bing图片链接" json:"url"`
	Hsh       string    `gorm:"size:32;comment:图片唯一hash_id" json:"hsh"`
	Desc      string    `gorm:"size:500;comment:备注描述" json:"desc"`
	ReleaseAt time.Time `gorm:"type:date;comment:发布日期" json:"releaseAt"`
	File      File      `json:"file"`
}
