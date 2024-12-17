package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/common/utils"
	"take-out/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommonController struct {
}

func (cc *CommonController) Upload(ctx *gin.Context) {
	code := e.SUCCESS
	// 获取前端传递的图片
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	// 拼接uuid的图片名称
	uuid := uuid.New()
	imageName := uuid.String() + file.Filename
	imagePath, err := utils.AliyunOss(imageName, file)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("AliyunOss upload failed", "err", err.Error())
	}
	ctx.JSON(http.StatusOK, common.Result{Code: code, Data: imagePath, Msg: e.GetMsg(code)})
}

// 图片上传
func (cc CommonController) LocalUpload(ctx *gin.Context) {
	code := e.SUCCESS
	file, err := ctx.FormFile("file")

	if err != nil {
		code = e.ERROR
		return
	}

	uuid := uuid.New()
	imageName := uuid.String() + file.Filename
	imagePath, err := utils.LocalFileSave(imageName, file)

	if err != nil {
		code = e.ERROR
		global.Log.Warn("Local upload failed", "err", err.Error())
	}
	ctx.JSON(http.StatusOK, common.Result{Code: code, Data: imagePath, Msg: e.GetMsg(code)})
}

// 图片上传访问
func (cc CommonController) LocalVisit(ctx *gin.Context) {

	// 从 URL 中获取子目录名和文件名
	subdir := ctx.Param("subdir")
	filename := ctx.Param("filename")

	filePath := utils.LocalFileVisit(subdir, filename)

	// 将文件内容发送给客户端
	ctx.File(filePath) // Gin 提供的便捷方法，用于从文件发送响应
}
