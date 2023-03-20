package userVo

type UserRegisterReqVo struct {
	Username    string `form:"username" json:"username" binding:"required,min=2"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6"`
}
