package userService

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type ResetPasswordReq struct {
	Password  string `json:"password"`
	Password1 string `json:"password1"`
}

func Login(c *gin.Context) *result.Result {
	r := result.New()
	db := c.Value("DB").(*gorm.DB)
	var params LoginReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		return r.Fail(err.Error())
	}

	splitHost := strings.Split(c.Request.Host, ":")
	if len(splitHost) < 1 {
		return r.Fail("host error")
	}
	var user model.User
	res := db.Where("username = ?", params.Username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return r.Fail("用户名或密码错误")
		}
		return r.Fail("")
	}

	//校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return r.Fail("用户名或密码错误")
	}

	// 生成jwtToken
	jwtConfig := config.Cfg.Jwt
	jwtToken, err := authUtil.GenerateJwtToken(jwtConfig, user.ID)
	if err != nil {
		return r.Fail("登录失败")
	}

	//生成散列hash token，并存到redis
	token, err := authUtil.SaveMd5Token(jwtToken)
	if err != nil {
		return r.Fail("登录失败")
	}
	c.SetCookie("token", token, jwtConfig.Expires, "/", splitHost[0], false, true)
	r.SetData(gin.H{"token": token})

	return r
}

func Register(c *gin.Context) *result.Result {
	r := result.New()
	var params RegisterReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		return r.FailErr(err)
	}
	if len(params.Username) < 4 {
		return r.Fail("用户名不能少于4位字符")
	}
	if len(params.Password) < 8 {
		return r.Fail("密码长度不能少于8位字符")
	}
	if !regexp.MustCompile("^1[345789]\\d{9}$").MatchString(params.PhoneNumber) {
		return r.Fail("请输入合法的手机号")
	}

	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	// 验证用户是否合法
	res := db.Model(&user).Where("username = ?", params.Username).Count(&count)
	if res.Error != nil {
		return r.FailErr(res.Error)
	}
	if count > 0 {
		return r.Fail("用户名不能重复")
	}
	res = db.Model(&user).Where("phone_number = ?", params.PhoneNumber).Count(&count)
	if res.Error != nil {
		return r.FailErr(res.Error)
	}
	if count > 0 {
		return r.Fail("手机号不能重复")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return r.FailType(common.UnknownError)
	}

	user = model.User{
		Username:    params.Username,
		Password:    string(hashPwd),
		PhoneNumber: params.PhoneNumber,
	}
	tx := db.Create(&user)
	if tx.Error != nil {
		return r.Fail("注册失败")
	}
	r.SetData(gin.H{"userId": user.ID})

	return r
}

func Logout(c *gin.Context) *result.Result {
	r := result.New()
	token := authUtil.GetToken(c)
	// 删除jwtToken
	if token != "" {
		config.RedisDb.Del(c, token)
	}
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
