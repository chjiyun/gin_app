package userService

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/userController/userVo"
	"gin_app/app/model"
	"gin_app/app/util/authUtil"
	"gin_app/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, reqVo userVo.UserLoginReqVo) (string, error) {
	db := c.Value("DB").(*gorm.DB)

	splitHost := strings.Split(c.Request.Host, ":")
	if len(splitHost) < 1 {
		return "", myError.New("host error")
	}
	var user model.User
	res := db.Where("name = ?", reqVo.Username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return "", myError.NewET(common.ErrUsernameOrPwd)
		}
		return "", res.Error
	}

	//校验密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqVo.Password))
	if err != nil {
		return "", myError.NewET(common.ErrUsernameOrPwd)
	}

	// 生成jwtToken
	jwtConfig := config.Cfg.Jwt
	jwtToken, err := authUtil.GenerateJwtToken(jwtConfig, user.ID)
	if err != nil {
		return "", myError.NewET(common.UnknownError)
	}

	//生成散列hash token，写入redis，token为键，jwtToken为值
	token, err := authUtil.SaveMd5Token(jwtToken)
	if err != nil {
		return "", myError.NewET(common.UnknownError)
	}
	c.SetCookie("token", token, jwtConfig.Expires, "/", splitHost[0], false, true)

	if config.Cfg.Env != gin.DebugMode {
		go saveLoginIpInfo(c.Copy(), token, user.ID)
	}

	return token, nil
}

func Register(c *gin.Context, reqVo userVo.UserRegisterReqVo) error {
	db := c.Value("DB").(*gorm.DB)
	var user model.User
	var count int64
	// 验证用户是否合法
	res := db.Model(&user).Where("name = ?", reqVo.Username).Count(&count)
	if res.Error != nil {
		return res.Error
	}
	if count > 0 {
		return myError.New("用户名不能重复")
	}
	res = db.Model(&user).Where("phone_number = ?", reqVo.PhoneNumber).Count(&count)
	if res.Error != nil {
		return res.Error
	}
	if count > 0 {
		return myError.New("手机号不能重复")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(reqVo.Password), bcrypt.DefaultCost)
	if err != nil {
		return myError.NewET(common.UnknownError)
	}

	user = model.User{
		Name:        reqVo.Username,
		Password:    string(hashPwd),
		PhoneNumber: reqVo.PhoneNumber,
	}
	tx := db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func Logout(c *gin.Context) {
	token := authUtil.GetToken(c)
	// 删除jwtToken
	if token != "" {
		config.RedisDb.Del(c, token)
	}
}
