package common

type PageRes struct {
	Count int64 `json:"count"`
	Rows  any   `json:"rows"`
}

type PageReq struct {
	Page     int `form:"page" json:"page" binding:"required,gt=0"`
	PageSize int `form:"pageSize" json:"pageSize" binding:"required,gt=0,lt=1000"`
}
