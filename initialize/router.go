package initialize

import (
	"take-out/internal/router"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

func routerInit() *gin.Engine {
	r := gin.Default()
	// 初始化websocket
	hub := service.NewHub()
	go hub.Run()
	return router.InitRouterGroup(r, hub)
}
