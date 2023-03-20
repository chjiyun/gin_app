package userService

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/model"
	"gin_app/app/util"
	"gin_app/app/util/authUtil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetCurrentUser 获取登录用户信息
func GetCurrentUser(c *gin.Context) *model.User {
	db := c.Value("DB").(*gorm.DB)
	userId := authUtil.GetSessionUserId(c)
	var user model.User
	db.Find(&user, userId)
	return &user
}

func GetUserPage(c *gin.Context, reqVo userVo.UserPageReqVo) (*common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var users []model.User
	var count int64
	// 初始条件可以放结构体里面
	tx := db.Model(model.User{})

	if reqVo.Keyword != "" {
		str := util.WriteString("%", reqVo.Keyword, "%")
		tx = tx.Where("name like ?", str).Or("phone_number like ?", str)
	}

	res := tx.Count(&count)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.Offset((reqVo.Page - 1) * reqVo.PageSize).Limit(reqVo.PageSize).Order("created_at").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return &common.PageRes{Count: count, Rows: users}, nil
}

// ResetPassword 重置密码  管理员才有权限
func ResetPassword(c *gin.Context, reqVo userVo.UserResetPasswordReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	userId := authUtil.GetSessionUserId(c)

	db.Model(&user).Where("id = ?", userId).Count(&count)
	if count == 0 {
		return false, errors.New("该用户不存在")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(reqVo.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, myError.NewET(common.UnknownError)
	}
	password := string(hashPwd)
	tx := db.Model(&user).Where("id = ?", userId).Update("password", password)
	if tx.Error != nil {
		return false, tx.Error
	}
	return true, nil
}
