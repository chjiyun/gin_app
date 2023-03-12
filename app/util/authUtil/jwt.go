package authUtil

import (
	"context"
	"crypto/md5"
	"fmt"
	"gin_app/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// CustomerClaims 嵌套匿名struct实现继承
type CustomerClaims struct {
	UserId uint `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(jwtConfig config.JwtConfig, userId uint) (string, error) {
	hmacSampleSecret := []byte(jwtConfig.SecretKey) //密钥，不能泄露
	token := jwt.New(jwt.SigningMethodHS256)
	nowTime := jwt.NewNumericDate(time.Now())
	expiredTime := jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.Expires) * time.Second))
	token.Claims = CustomerClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: nowTime,          // 签名生效时间
			ExpiresAt: expiredTime,      // 签名过期时间
			Issuer:    jwtConfig.Issuer, // 签名颁发者
		},
	}
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func ParseJwtToken(tokenString string, secret string) (*CustomerClaims, error) {
	var hmacSampleSecret = []byte(secret)
	// 前面例子生成的token
	token, err := jwt.ParseWithClaims(tokenString, &CustomerClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})

	if err != nil {
		return nil, err
	}
	claims := token.Claims.(*CustomerClaims)
	return claims, nil
}

func SaveMd5Token(jwtToken string) (string, error) {
	//生成散列hash 并存到redis
	hash := md5.Sum([]byte(jwtToken))
	token := fmt.Sprintf("%x", hash)
	expiration := time.Duration(config.Cfg.Jwt.Expires) * time.Second
	_, err := config.RedisDb.Set(context.Background(), token, jwtToken, expiration).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}
