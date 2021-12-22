package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// soft_delete 此插件还支持混合模式，删除同时更新删除时间
// 文件model
type File struct {
	ID        uint                  `gorm:"primaryKey" json:"id"`
	Uid       string                `gorm:"size:36;not null;comment:唯一id" json:"uid"`
	Ext       string                `gorm:"size:20;comment:文件后缀" json:"ext"`
	Type      string                `gorm:"size:20;comment:文件类型: image,video,audio,txt" json:"type"`
	MimeType  string                `gorm:"size:100" json:"mimeType"`
	Name      string                `gorm:"size:100;comment:原始文件名" json:"name"`
	LocalName string                `gorm:"size:39;comment:服务器本地文件名（取雪花id）" json:"localName"`
	Path      string                `gorm:"size:100;comment:文件存储相对路径" json:"path"`
	Size      uint                  `gorm:"comment:文件大小（Byte）" json:"size"`
	Desc      string                `gorm:"size:255;comment:文件描述" json:"desc"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag;not null;default:0" json:"-"`
}
