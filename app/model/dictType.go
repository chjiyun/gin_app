package model

import "gin_app/app/common"

type DictType struct {
	common.BaseModel
	Name      string      `gorm:"size:100;comment:名称" json:"name"`
	Value     string      `gorm:"size:100;comment:唯一标识符" json:"value"`
	DictValue []DictValue `gorm:"foreignKey: TypeId" json:"dictValues"`
}
