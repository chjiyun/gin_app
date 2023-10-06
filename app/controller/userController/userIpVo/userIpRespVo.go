package userIpVo

import (
	"gin_app/app/common"
	"gin_app/app/model"
	"time"
)

type UserIpRespVo struct {
	common.BaseModel
	UserId    uint64      `gorm:"not null;comment:用户id" json:"userId"`
	Ip        string      `gorm:"size:45;comment:ip地址" json:"ip"`
	Country   string      `gorm:"size:30;comment:国家" json:"country"`
	Province  string      `gorm:"size:30;comment:省份" json:"province"`
	City      string      `gorm:"size:30;comment:城市" json:"city"`
	Area      string      `gorm:"size:30;comment:市、区" json:"area"`
	CreatedAt time.Time   `json:"createdAt"`
	User      *model.User `json:"user"`
}

func (UserIpRespVo) TableName() string {
	return "user_ip"
}
