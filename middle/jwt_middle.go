package middle

import (
	"net/http"
	"strconv"
	"take-out/common"
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/common/utils"
	"take-out/global"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// VerifyJWTAdmin 验证管理员Token
func VerifyJWTRedis(secret string, UserId uint64) bool {
	_, err := global.Redis.Get(secret + strconv.FormatUint(UserId, 10)).Result()
	if err == redis.Nil || err != nil {
		return true
	}
	return false
}

func VerifyJWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		admin := global.Config.Jwt.Admin
		token := c.Request.Header.Get(admin.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, admin.Secret)
		isOk := VerifyJWTRedis(admin.Secret, payLoad.UserId)
		if err != nil || isOk {
			code = e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.GrantScope)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}

func VerifyJWTUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		token := c.Request.Header.Get(global.Config.Jwt.User.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, global.Config.Jwt.User.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.GrantScope)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}
