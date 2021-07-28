package main

import (
	"RCM_CS/APIs"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	//r.POST("/uploadfile", APIs.UploadHandler)

	r.POST("/upload", APIs.UploadByJson)

	//r.GET("/downloadfile/:filename", APIs.DownloadHandler)

	r.GET("/download", APIs.Download)

	r.GET("/pull", APIs.Pull)

	r.PUT("/update", APIs.UpdateField)

	return r
}
