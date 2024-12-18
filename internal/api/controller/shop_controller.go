package controller

import (
	"log/slog"
	"net/http"
	"strconv"
	"take-out/common"
	"take-out/common/e"
	"take-out/global"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Shop struct{}

func (s Shop) GetShopStatus(c *gin.Context) {
	code := e.SUCCESS
	// 从redis中获取状态
	shopStatus, err := global.Redis.Get("shopStatus").Result()
	// redis中没有赋默认值
	if err == redis.Nil {
		shopStatus = "1"
	} else if err != nil {
		code = e.ERROR
		slog.Warn("获取营业状态失败", "Err:", err.Error())
	}
	// 将状态转为int类型
	status, _ := strconv.Atoi(shopStatus)
	c.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: status,
		Msg:  e.GetMsg(code),
	})
}

func (s Shop) SetShopStatus(c *gin.Context) {
	code := e.SUCCESS
	status := c.Param("status")
	// 设置状态， 过期时间为永久
	err := global.Redis.Set("shopStatus", status, 0).Err()

	if err != nil {
		code = e.ERROR
		slog.Warn("设置营业状态失败", "Err:", err.Error())
	}

	c.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
