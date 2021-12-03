package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// soft_delete 此插件还支持混合模式，删除同时更新删除时间
// 文件model
type File struct {
	ID        uint                  `gorm:"primaryKey" json:"id"`
	Type      string                `gorm:"size:20;comment:文件类型" json:"type"`
	MimeType  string                `gorm:"size:100" json:"mimeType"`
	Name      string                `gorm:"size:100;comment:原文件名" json:"name"`
	Path      string                `gorm:"size:250;comment:存储路径" json:"path"`
	Size      uint                  `gorm:"comment:文件大小（Byte）" json:"size"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}
