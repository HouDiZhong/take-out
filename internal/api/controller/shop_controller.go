package controller

import (
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
	shopStatus, err := global.Redis.Get("shopStatus").Result()

	if err == redis.Nil {
		shopStatus = "1"
	} else if err != nil {
		code = e.ERROR
	}

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
	err := global.Redis.Set("shopStatus", status, 0).Err()

	if err != nil {
		code = e.ERROR
	}

	c.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
