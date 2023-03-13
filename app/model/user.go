package model

// User 用户
type User struct {
	BaseModel
	Username    string `gorm:"size:100;comment:用户名" json:"username"`
	PhoneNumber string `gorm:"size:20;comment:手机号" json:"phoneNumber"`
	Email       string `gorm:"size:100;comment:邮箱" json:"email"`
	Password    string `gorm:"size:80;comment:密码" json:"password"`
	Gender      int    `gorm:"comment:性别" json:"gender"`
	Portrait    uint64 `gorm:"comment:头像" json:"portrait"`
}
