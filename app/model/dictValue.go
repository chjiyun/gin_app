package model

import "gin_app/app/common"

// DictValue 字典值映射
type DictValue struct {
	common.BaseModel
	Label  string `gorm:"size:100;comment:字典值对应中文名" json:"label"`
	Value  string `gorm:"not null;size:100;comment:字典值" json:"value"`
	Sort   int    `gorm:"comment:序号" json:"sort"`
	Enable bool   `gorm:"not null;default:true;comment:是否启用" json:"enable"`
	TypeId uint64 `gorm:"not null;comment:类型id" json:"type_id"`
}
