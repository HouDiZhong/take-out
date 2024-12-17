package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"take-out/global"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func AliyunOss(fileName string, file *multipart.FileHeader) (string, error) {
	config := global.Config.AliOss
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return "", err
	}

	fileData, _ := file.Open()
	defer fileData.Close()

	err = bucket.PutObject(fileName, fileData)
	if err != nil {
		return "", err
	}
	imagePath := "https://" + config.BucketName + "." + config.EndPoint + "/" + fileName
	fmt.Println("文件上传到：", imagePath)
	return imagePath, nil
}

func LocalFileSave(fileName string, file *multipart.FileHeader) (string, error) {
	config := global.Config.LocalPath

	now := time.Now()
	ext := path.Ext(fileName)
	fileName = strconv.Itoa(now.Nanosecond()) + ext
	filePath := fmt.Sprintf("%s%s%s%s",
		config.UploadDir,
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%04d", now.Day()))

	CreateDir(filePath)
	fullPath := filePath + "/" + fileName
	err := SaveUploadedFile(file, fullPath)
	if err != nil {
		return "", err
	}
	// nginx 反代将/admin 代理成了/api
	fullPath = config.ImageHost + "api/" + fullPath
	return fullPath, nil
}

func LocalFileVisit(subdir, filename string) string {
	config := global.Config.LocalPath
	// 构建图片的完整路径
	filePath := filepath.Join(config.UploadDir, subdir, filename)

	return filePath
}

// CreateDir 创建目录
func CreateDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// IsExist 判断是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
