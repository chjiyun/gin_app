package fileVo

type FileRespVo struct {
	ID   string `gorm:"primaryKey;size:32;comment:nanoid拼接ext" json:"id"`
	Uid  uint64 `gorm:"not null;comment:雪花id" json:"uid"`
	Ext  string `gorm:"size:20;comment:文件后缀" json:"ext"`
	Type string `gorm:"size:20;comment:文件类型: image,video,audio,txt" json:"type"`
	Name string `gorm:"size:100;comment:原始文件名" json:"name"`
	Size uint   `gorm:"comment:文件大小（Byte）" json:"size"`
	Desc string `gorm:"size:255;comment:文件描述" json:"desc"`
}
