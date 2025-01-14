package initialize

import (
	"take-out/config"
	"take-out/global"
	"take-out/logger"

	"github.com/gin-gonic/gin"
)

func GlobalInit() *gin.Engine {
	// 配置文件初始化
	global.Config = config.InitLoadConfig()
	// Log初始化
	global.Log = logger.NewMySlog(global.Config.Log.Level, global.Config.Log.FilePath)
	// Gorm初始化
	global.DB = InitDatabase(global.Config.DataSource.Dsn())
	// Redis初始化
	global.Redis = initRedis()
	// Router初始化
	router := routerInit()
	// Cron初始化
	initCron()
	// WebSocket初始化
	// InitWebSocket(router)
	return router
}
