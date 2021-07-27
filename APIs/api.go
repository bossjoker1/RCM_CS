package APIs

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var dir, _ = os.Getwd()

func UploadByJson(c *gin.Context) {

	data := make(map[string]interface{})
	data["uid"] = ""
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	for k, v := range data {
		fmt.Println("key: ", k, "value: ", v)
	}
	jsonByte, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal failed. %v\n")
		return
	}
	//fmt.Printf("%s\n", jsonByte)

	// 根据json byte数据创建本地json文件
	// 创建的文件名根据用户发来的uid组成
	dst := fmt.Sprintf(".\\files\\" + data["uid"].(string) + ".json")

	err = ioutil.WriteFile(dst, jsonByte, 0644)

	if err != nil {
		log.Printf("write file failed. %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "successfully upload",
	})
}

//
//func UploadHandler(c *gin.Context) {
//	file, err := c.FormFile("upload")
//
//	if err != nil {
//		log.Printf("read formfile failed. %v\n", err)
//
//		c.JSON(http.StatusInternalServerError,
//			gin.H{
//				"code": 400,
//				"msg":  fmt.Sprintf("read formfile failed. %v\n", err),
//			})
//		return
//	}
//
//	dst := fmt.Sprintf(".\\files\\" + file.Filename)
//	fmt.Println("env :", dir)
//
//	// 保存到指定路径
//	err = c.SaveUploadedFile(file, dst)
//
//	if err != nil {
//		log.Printf("save the file failed. %v\n", err)
//		c.JSON(http.StatusInternalServerError,
//			gin.H{
//				"code": 500,
//				"msg":  fmt.Sprintf("save the file failed. %v\n", err),
//			})
//
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"msg":      200,
//		"filepath": dst,
//	})
//}

// 下载文件，返回的是
func Download(c *gin.Context) {
	req := make(map[string]interface{})

	req["uid"] = ""

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	// 打开文件

	fileBytes, err := ioutil.ReadFile(".\\files\\" + req["uid"].(string) + ".json")
	if err != nil {
		log.Printf("read the file failed. %v\n", err)
		return
	}
	// 将json文件以string返回
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": fmt.Sprintf("%s", fileBytes),
	})

}

//func DownloadHandler(c *gin.Context) {
//	filename := c.Param("filename")
//
//	_, err := os.Open(".\\files\\" + filename)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError,
//			gin.H{
//				"msg": " file not existed.",
//			})
//		return
//	}
//
//	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
//	c.Writer.Header().Add("Content-Type", "application/octet-stream")
//
//	// 浏览器下载文件
//	c.File(".\\files\\" + filename)
//
//	//c.Data()
//	//c.JSON(http.StatusInternalServerError,
//	//	gin.H{
//	//		"code" : 200,
//	//		"msg"  : "successfully pull",
//	//	})
//	return
//
//}
