package user

import (
	"errors"
	"gin_app/app/model"
	"gin_app/app/result"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"

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

func Login(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	var params LoginReq
	r := result.New()

	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail(err.Error()))
		return
	}

	splitHost := strings.Split(c.Request.Host, ":")
	if len(splitHost) < 1 {
		c.JSON(http.StatusOK, r.Fail("host error"))
		return
	}
	var user model.User
	res := db.Select("id").Where("username = ?", params.Username).First(&user)
	if res.Error != nil {
		r.Fail("")
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			r.Fail("用户名或密码错误")
		}
		c.JSON(http.StatusOK, r)
		return
	}

	//校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		c.JSON(http.StatusOK, r.Fail("用户名或密码错误"))
		return
	}

	// 生成jwtToken
	jwtConfig := config.Cfg.Jwt
	jwtToken, err := authUtil.GenerateJwtToken(jwtConfig, user.ID)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail("登录失败"))
		return
	}

	//生成散列hash token，并存到redis
	token, err := authUtil.SaveMd5Token(jwtToken)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail("登录失败"))
		return
	}

	c.SetCookie("token", token, jwtConfig.Expires, "/", splitHost[0], false, true)
	r.SetData(gin.H{"token": token})

	c.JSON(http.StatusOK, r)
}

func Register(c *gin.Context) {
	var params RegisterReq
	r := result.New()

	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail(err.Error()))
		return
	}
	db := c.Value("DB").(*gorm.DB)
	var user model.User
	// 验证用户是否合法
	res := db.Where("username = ?", params.Username).First(&user)
	if res.Error != nil {
		r.Fail("")
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			r.Fail("用户名不能重复")
		}
		c.JSON(http.StatusOK, r)
		return
	}
	res = db.Where("phone_number = ?", params.PhoneNumber).First(&user)
	if res.Error != nil {
		r.Fail("")
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			r.Fail("手机号不能重复")
		}
		c.JSON(http.StatusOK, r)
		return
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, r.Fail(err.Error()))
		return
	}

	user = model.User{
		Username:    params.Username,
		Password:    string(hashPwd),
		PhoneNumber: params.PhoneNumber,
	}
	tx := db.Create(&user)
	if tx.Error != nil {
		c.JSON(http.StatusOK, r.Fail("注册失败"))
		return
	}

	r.SetData(gin.H{"userId": user.ID})
	c.JSON(http.StatusOK, r)
}

func Logout(c *gin.Context) {
	token := authUtil.GetToken(c)
	// 删除jwtToken
	config.RedisDb.Del(c, token)
}
