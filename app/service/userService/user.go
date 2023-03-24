package userService

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/userController/userIpVo"
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/model"
	"gin_app/app/service/cacheService"
	"gin_app/app/service/toolService"
	"gin_app/app/util"
	"gin_app/app/util/authUtil"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetCurrentUser 获取登录用户信息
func GetCurrentUser(c *gin.Context) *userVo.UserRespVo {
	db := c.Value("DB").(*gorm.DB)
	userId := authUtil.GetSessionUserId(c)

	var user model.User
	db.Find(&user, userId)
	var userIp model.UserIp
	token := authUtil.GetToken(c)
	if userIpId, err := cacheService.GetSessionIp(token); err == nil {
		db.Find(&userIp, userIpId)
	}

	userRespVo := userVo.UserRespVo{}
	_ = copier.Copy(&userRespVo, &user)
	_ = copier.Copy(userRespVo.UserIp, &userIp)
	return &userRespVo
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

func saveLoginIpInfo(c *gin.Context, token string, userId uint) {
	log := c.Value("Logger").(*logrus.Entry)
	db := c.Value("DB").(*gorm.DB)

	ip := c.ClientIP()
	info, err := toolService.GetIpInfo(c, ip)
	if err != nil {
		log.Error(err)
		return
	}

	var userIp = model.UserIp{
		UserId: userId,
	}
	if err = copier.Copy(&userIp, &info); err != nil {
		log.Error(err)
		return
	}
	db.Create(&userIp)

	if err = cacheService.SaveSessionIp(token, userIp.ID); err != nil {
		log.Error(err)
	}
}

func GetUserIpPage(c *gin.Context, reqVo userIpVo.UserIpPageReqVo) (*common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var userIps []userIpVo.UserIpRespVo
	var count int64

	db.Joins("user").Model(&userIps).Count(&count)
	tx := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("name")
	})
	if reqVo.StartTime != nil {
		tx = tx.Where("created_at >= ?", reqVo.StartTime)
	}
	if reqVo.EndTime != nil {
		tx = tx.Where("created_at < ?", reqVo.EndTime)
	}

	tx.Offset((reqVo.Page - 1) * reqVo.PageSize).
		Limit(reqVo.PageSize).
		Order("created_at desc").
		Find(&userIps)

	return &common.PageRes{Count: count, Rows: userIps}, nil
}
