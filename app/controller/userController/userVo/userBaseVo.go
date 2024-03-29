package userVo

type UserBaseVo struct {
	ID          uint64 `json:"id"`
	Name        string `gorm:"size:100;comment:用户名" json:"name"`
	PhoneNumber string `gorm:"size:20;comment:手机号" json:"phoneNumber"`
	Email       string `gorm:"size:100;comment:邮箱" json:"email"`
	Gender      int    `gorm:"comment:性别" json:"gender"`
	Portrait    string `gorm:"comment:头像" json:"portrait"`
}
