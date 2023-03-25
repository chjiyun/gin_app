package common

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	dateTime := fmt.Sprintf("%q", time.Time(d).Format("2006-01-02"))
	return []byte(dateTime), nil
}

// BaseModel soft_delete 此插件还支持混合模式，删除同时更新删除时间
type BaseModel struct {
	ID        uint                  `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"-"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag;not null;default:0" json:"-"`
}
