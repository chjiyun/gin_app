package model

// 文件model
type File struct {
	BaseModel
	Uid       uint64  `gorm:"not null;comment:雪花id" json:"uid"`
	Ext       string  `gorm:"size:20;comment:文件后缀" json:"ext"`
	Type      string  `gorm:"size:20;comment:文件类型: image,video,audio,txt" json:"type"`
	MimeType  string  `gorm:"size:100" json:"mimeType"`
	Name      string  `gorm:"size:100;comment:原始文件名" json:"name"`
	LocalName string  `gorm:"size:39;comment:服务器本地文件名（取雪花id）" json:"-"`
	Path      string  `gorm:"size:100;comment:文件存储相对路径" json:"-"`
	Size      uint    `gorm:"comment:文件大小（Byte）" json:"size"`
	Desc      string  `gorm:"size:255;comment:文件描述" json:"desc"`
	Thumbs    []Thumb `json:"thumbs"`
}
