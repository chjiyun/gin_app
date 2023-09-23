package model

import "gin_app/app/common"

// DictValue 字典值映射
type DictValue struct {
	common.BaseModel
	Label  string `gorm:"size:100;comment:字典值对应中文名" json:"label"`
	Value  int16  `gorm:"not null;comment:字典值" json:"value"`
	Sort   int32  `gorm:"comment:序号" json:"sort"`
	Enable bool   `gorm:"not null;default:true;comment:是否启用" json:"enable"`
	TypeId int64  `gorm:"not null;comment:类型id" json:"type_id"`
}
