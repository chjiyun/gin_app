package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Bing struct {
	ID        uint                  `gorm:"primaryKey" json:"id"`
	FileId    uint                  `gorm:"not null;comment:外键" json:"file_id"`
	Url       string                `gorm:"size:255;comment:原bing图片链接" json:"url"`
	Desc      string                `gorm:"size:500;comment:备注描述" json:"desc"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag;not null;default:0" json:"-"`
}
