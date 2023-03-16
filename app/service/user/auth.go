package user

import (
	"errors"
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

func Login(c *gin.Context, r *result.Result) {
	db := c.Value("DB").(*gorm.DB)
	var params LoginReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		r.Fail(err.Error())
		return
	}

	splitHost := strings.Split(c.Request.Host, ":")
	if len(splitHost) < 1 {
		r.Fail("host error")
		return
	}
	var user model.User
	res := db.Where("username = ?", params.Username).First(&user)
	if res.Error != nil {
		r.Fail("")
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			r.Fail("用户名或密码错误")
		}
		return
	}

	//校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		r.Fail("用户名或密码错误")
		return
	}

	// 生成jwtToken
	jwtConfig := config.Cfg.Jwt
	jwtToken, err := authUtil.GenerateJwtToken(jwtConfig, user.ID)
	if err != nil {
		r.Fail("登录失败")
		return
	}

	//生成散列hash token，并存到redis
	token, err := authUtil.SaveMd5Token(jwtToken)
	if err != nil {
		r.Fail("登录失败")
		return
	}

	c.SetCookie("token", token, jwtConfig.Expires, "/", splitHost[0], false, true)
	r.SetData(gin.H{"token": token})

}

func Register(c *gin.Context, r *result.Result) {
	var params RegisterReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		r.Fail(err.Error())
		return
	}
	if len(params.Username) < 4 {
		r.Fail("用户名不能少于4位字符")
		return
	}
	if len(params.Password) < 8 {
		r.Fail("密码长度不能少于8位字符")
		return
	}
	if !regexp.MustCompile("^1[345789]\\d{9}$").MatchString(params.PhoneNumber) {
		r.Fail("请输入合法的手机号")
		return
	}

	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	// 验证用户是否合法
	res := db.Model(&user).Where("username = ?", params.Username).Count(&count)
	if res.Error != nil {
		r.Fail(res.Error.Error())
		return
	}
	if count > 0 {
		r.Fail("用户名不能重复")
		return
	}
	res = db.Model(&user).Where("phone_number = ?", params.PhoneNumber).Count(&count)
	if res.Error != nil {
		r.Fail(res.Error.Error())
		return
	}
	if count > 0 {
		r.Fail("手机号不能重复")
		return
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		r.Fail(err.Error())
		return
	}

	user = model.User{
		Username:    params.Username,
		Password:    string(hashPwd),
		PhoneNumber: params.PhoneNumber,
	}
	tx := db.Create(&user)
	if tx.Error != nil {
		r.Fail("注册失败")
		return
	}

	r.SetData(gin.H{"userId": user.ID})
}

func Logout(c *gin.Context) {
	token := authUtil.GetToken(c)
	// 删除jwtToken
	if token != "" {
		config.RedisDb.Del(c, token)
	}
}

// ResetPassword 重置密码  管理员才有权限
func ResetPassword(c *gin.Context, r *result.Result) {
	var params ResetPasswordReq

	err := c.ShouldBindJSON(&params)
	if err != nil {
		r.Fail(err.Error())
		return
	}
	if params.Password != params.Password1 {
		r.Fail("密码不一致")
		return
	}

	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	userId := authUtil.GetSessionUserId(c)

	db.Model(&user).Where("id = ?", userId).Count(&count)
	if count == 0 {
		r.Fail("该用户不存在")
		return
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		r.Fail(err.Error())
		return
	}
	password := string(hashPwd)
	tx := db.Model(&user).Where("id = ?", userId).Update("password", password)
	if tx.Error != nil {
		r.Fail("更新失败")
		return
	}

}
