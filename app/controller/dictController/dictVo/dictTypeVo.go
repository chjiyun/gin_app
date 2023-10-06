package dictVo

import "gin_app/app/common"

type DictTypeCreateReqVo struct {
	Name  string `gorm:"size:100;comment:名称" json:"name" binding:"required"`
	Value string `gorm:"size:100;comment:唯一标识符" json:"value" binding:"required"`
}

type DictTypeUpdateReqVo struct {
	ID    uint64 `json:"id" binding:"required"`
	Name  string `gorm:"size:100;comment:名称" json:"name"`
	Value string `gorm:"size:100;comment:唯一标识符" json:"value"`
}

type DictTypeRespVo struct {
	common.BaseModel
	Name      string            `gorm:"size:100;comment:名称" json:"name"`
	Value     string            `gorm:"size:100;comment:唯一标识符" json:"value"`
	DictValue []DictValueRespVo `json:"dictValues"`
}
