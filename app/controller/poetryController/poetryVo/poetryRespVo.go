package poetryVo

type PoetryRespVo struct {
	ID            uint64                `json:"id"`
	Author        string                `gorm:"size:50;comment:作者" json:"author"`
	Title         string                `gorm:"not null;size:255;comment:作品名" json:"title"`
	Desc          string                `gorm:"size:1000;comment:解释" json:"desc"`
	Tag           int16                 `gorm:"comment:标签" json:"tag"`
	PoetryContent []PoetryContentRespVo `json:"poetryContent"`
}
