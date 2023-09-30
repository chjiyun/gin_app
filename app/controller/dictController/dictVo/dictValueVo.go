package dictVo

import "gin_app/app/common"

type DictValueCreateReqVo struct {
	Label  string `gorm:"size:100;comment:字典值对应中文名" json:"label" binding:"required"`
	Value  string `gorm:"not null;comment:字典值" json:"value" binding:"required"`
	Sort   int32  `gorm:"comment:序号" json:"sort"`
	Enable bool   `gorm:"not null;default:true;comment:是否启用" json:"enable"`
	TypeId uint64 `gorm:"not null;comment:类型id" json:"type_id" binding:"required"`
}

type DictValueUpdateReqVo struct {
	ID     uint64 `json:"id" binding:"required"`
	Label  string `gorm:"size:100;comment:字典值对应中文名" json:"label"`
	Value  string `gorm:"not null;comment:字典值" json:"value"`
	Sort   int32  `gorm:"comment:序号" json:"sort"`
	Enable bool   `gorm:"not null;default:true;comment:是否启用" json:"enable"`
}

type DictValueRespVo struct {
	common.BaseModel
	Label  string `gorm:"size:100;comment:字典值对应中文名" json:"label"`
	Value  string `gorm:"not null;comment:字典值" json:"value"`
	Sort   int32  `gorm:"comment:序号" json:"sort"`
	Enable bool   `gorm:"not null;default:true;comment:是否启用" json:"enable"`
	TypeId uint64 `gorm:"not null;comment:类型id" json:"type_id"`
}

type DictValueReqVo struct {
	Keyword   string `form:"keyword" json:"keyword"`
	TypeId    uint64 `form:"type_id" json:"type_id" binding:"required"`
	SortField string `form:"sortField" json:"sortField"`
	SortOrder string `form:"SortOrder" json:"SortOrder"`
}
