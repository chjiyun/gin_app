package model

// 缩略图 webp/jpeg/jpg/png model
type Thumb struct {
	BaseModel
	Uid       uint64 `gorm:"not null;comment:雪花id" json:"uid"`
	FileId    uint   `gorm:"not null;comment:file_id" json:"file_id"`
	Ext       string `gorm:"size:20;comment:文件后缀" json:"ext"`
	Width     uint   `gorm:"comment:宽度" json:"width"`
	Height    uint   `gorm:"comment:高度" json:"height"`
	Name      string `gorm:"size:120;comment:原始文件名" json:"name"`
	LocalName string `gorm:"size:39;comment:服务器本地文件名（取雪花id）" json:"localName"`
	Path      string `gorm:"size:100;comment:文件存储相对路径" json:"path"`
	Size      uint   `gorm:"comment:文件大小（Byte）" json:"size"`
}
