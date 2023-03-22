package model

import "gin_app/app/common"

type UserIp struct {
	common.BaseModel
	UserId   uint   `gorm:"comment:用户名" json:"userId"`
	Ip       string `gorm:"size:45;comment:ip地址" json:"ip"`
	Country  string `gorm:"size:30;comment:国家" json:"country"`
	Province string `gorm:"size:30;comment:省份" json:"province"`
	City     string `gorm:"size:30;comment:城市" json:"city"`
	Area     string `gorm:"size:30;comment:市、区" json:"area"`
}
