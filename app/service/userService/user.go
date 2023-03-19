package userService

import (
	"fmt"
	"gin_app/app/common"
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetCurrentUser(c *gin.Context) {
	r := c.Value("Result").(*result.Result)
	user := getSessionUser(c)
	r.SetData(user)
}

// getSessionUser 获取登录用户信息
func getSessionUser(c *gin.Context) model.User {
	db := c.Value("DB").(*gorm.DB)
	userId := authUtil.GetSessionUserId(c)
	var user model.User
	db.Find(&user, userId)
	return user
}

func GetPageUsers(c *gin.Context, reqVo userVo.UserPageReqVo) {
	r := c.Value("Result").(*result.Result)
	db := c.Value("DB").(*gorm.DB)

	var users []model.User
	var count int64
	// 初始条件可以放结构体里面
	tx := db.Model(model.User{})

	if reqVo.Keyword != "" {
		str := fmt.Sprintf("%%%s%%", reqVo.Keyword)
		tx = tx.Where("name like ?", str).Or("phone_number like ?", str)
	}

	res := tx.Count(&count)
	if res.Error != nil {
		r.FailErr(res.Error)
		return
	}

	res = tx.Offset((reqVo.Page - 1) * reqVo.PageSize).Limit(reqVo.PageSize).Order("created_at").Find(&users)
	if res.Error != nil {
		r.FailErr(res.Error)
		return
	}

	r.SetData(common.PageRes{
		Count: count,
		Rows:  users,
	})
}

// ResetPassword 重置密码  管理员才有权限
func ResetPassword(c *gin.Context) *result.Result {
	r := result.New()
	var params ResetPasswordReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		return r.FailErr(err)
	}
	if params.Password != params.Password1 {
		return r.Fail("密码不一致")
	}

	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	userId := authUtil.GetSessionUserId(c)

	db.Model(&user).Where("id = ?", userId).Count(&count)
	if count == 0 {
		return r.Fail("该用户不存在")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return r.FailType(common.UnknownError)
	}
	password := string(hashPwd)
	tx := db.Model(&user).Where("id = ?", userId).Update("password", password)
	if tx.Error != nil {
		return r.Fail("")
	}
	return r
}
