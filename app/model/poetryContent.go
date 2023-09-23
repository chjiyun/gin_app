package model

import "gin_app/app/common"

type PoetryContent struct {
	common.BaseModel
	PoetryId uint64 `gorm:"not null;comment:作者" json:"poetry_id"`
	Content  string `gorm:"comment:作品名" json:"content"`
	Sort     int    `gorm:"comment:标签" json:"sort"`
}

func (PoetryContent) TableName() string {
	return "b_poetry_content"
}
