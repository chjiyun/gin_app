package dictVo

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
	ID         uint64            `json:"id"`
	Name       string            `gorm:"size:100;comment:名称" json:"name"`
	Value      string            `gorm:"size:100;comment:唯一标识符" json:"value"`
	DictValues []DictValueRespVo `json:"dictValues"`
}
