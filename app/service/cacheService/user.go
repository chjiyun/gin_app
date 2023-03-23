package cacheService

import (
	"context"
	"gin_app/config"
	"time"
)

func SaveSessionIp(token string, userIpId uint) error {
	key := "login_ip:" + token
	expiration := time.Duration(config.Cfg.Jwt.Expires) * time.Second
	_, err := config.RedisDb.Set(context.Background(), key, userIpId, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetSessionIp(token string) (uint, error) {
	key := "login_ip:" + token
	userIpId, err := config.RedisDb.Get(context.Background(), key).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(userIpId), nil
}
