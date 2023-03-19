package userService

import (
	"fmt"
	"gin_app/app/common"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

func GetCurrentUser(c *gin.Context) *result.Result {
	r := result.New()
	user := getSessionUser(c)
	r.SetData(user)
	return r
}

// getSessionUser 获取登录用户信息
func getSessionUser(c *gin.Context) model.User {
	db := c.Value("DB").(*gorm.DB)
	userId := authUtil.GetSessionUserId(c)
	var user model.User
	db.Find(&user, userId)
	return user
}

func GetPageUsers(c *gin.Context) *result.Result {
	r := result.New()
	db := c.Value("DB").(*gorm.DB)

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	keyword := c.Query("keyword")
	if page < 1 || pageSize < 1 {
		return r.Fail("参数缺失")
	}

	var users []model.User
	var count int64
	// 初始条件可以放结构体里面
	tx := db.Model(model.User{})

	if keyword != "" {
		str := fmt.Sprintf("%%%s%%", keyword)
		tx = tx.Where("name like ?", str).Or("phone_number like ?", str)
	}

	res := tx.Count(&count)
	if res.Error != nil {
		return r.FailErr(res.Error)
	}

	res = tx.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at").Find(&users)
	if res.Error != nil {
		return r.FailErr(res.Error)
	}

	r.SetData(common.Page{
		Count: count,
		Rows:  users,
	})
	return r
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
