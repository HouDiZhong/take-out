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
		admin := global.Config.Jwt.Admin
		token := c.Request.Header.Get(admin.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, admin.Secret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, common.Result{Code: e.UNKNOW_IDENTITY})
			c.Abort()
			return
		}

		if isOk := VerifyJWTRedis(admin.Secret, payLoad.UserId); isOk {
			c.JSON(http.StatusUnauthorized, common.Result{Code: e.UNKNOW_IDENTITY})
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
		user := global.Config.Jwt.User
		token := c.Request.Header.Get(user.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, user.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.Result{Code: e.UNKNOW_IDENTITY})
			c.Abort()
			return
		}

		if isOk := VerifyJWTRedis(user.Secret, payLoad.UserId); isOk {
			c.JSON(http.StatusUnauthorized, common.Result{Code: e.UNKNOW_IDENTITY})
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
