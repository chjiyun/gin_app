package userVo

// UserResetPasswordReqVo 用户密码重置请求Vo
type UserResetPasswordReqVo struct {
	Password  string `form:"password" json:"password" binding:"required"`
	Password1 string `form:"password1" json:"password1" binding:"required,eqfield=password"`
}
