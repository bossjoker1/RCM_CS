package main

import (
	"RCM_CS/APIs"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	// json数据上传后服务端写入创建文件
	// 直接整个文件的上传和下载形式

	r.POST("/upload", APIs.UploadByProperties)

	r.GET("/downloadfile/:filename", APIs.DownloadHandler)

	r.GET("/pull", APIs.PersonalizedPull)

	r.PUT("/update", APIs.PersonalizedUpdate)

	return r
}
