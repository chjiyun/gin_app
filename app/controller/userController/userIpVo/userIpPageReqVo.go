package userIpVo

import (
	"gin_app/app/common"
	"time"
)

type UserIpPageReqVo struct {
	common.PageReq
	Ip        string    `form:"ip" json:"ip"`
	UserId    uint64    `form:"userId" json:"userId"`
	StartTime time.Time `form:"startTime" json:"startTime"`
	EndTime   time.Time `form:"endTime" json:"endTime"`
}
