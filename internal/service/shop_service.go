package service

import (
	"log/slog"
	"strconv"
	"take-out/common/e"
	"take-out/global"

	"github.com/go-redis/redis"
)

type ShopService interface {
	GetShopService(code *int) (status int)
	SetShopService(status string, code *int)
}

type ShopServiceImpl struct{}

func NewShopService() ShopService {
	return ShopServiceImpl{}
}

func (s ShopServiceImpl) GetShopService(code *int) (status int) {
	// 从redis中获取状态
	shopStatus, err := global.Redis.Get("shopStatus").Result()
	// redis中没有赋默认值
	if err == redis.Nil {
		shopStatus = "1"
	} else if err != nil {
		*code = e.ERROR
		slog.Warn("获取营业状态失败", "Err:", err.Error())
	}
	// 将状态转为int类型
	status, _ = strconv.Atoi(shopStatus)
	return
}

func (s ShopServiceImpl) SetShopService(status string, code *int) {
	// 设置状态， 过期时间为永久
	err := global.Redis.Set("shopStatus", status, 0).Err()
	if err != nil {
		*code = e.ERROR
		slog.Warn("设置营业状态失败", "Err:", err.Error())
	}
}
