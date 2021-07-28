package APIs

import (
	"RCM_CS/Models"
	"RCM_CS/Utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

//var dir, _ = os.Getwd()

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
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "we can not find the file, please make sure the path is right.",
		})
		return
	}
	// 将json文件以string返回
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": fmt.Sprintf("%s", fileBytes),
	})

}

// 配置文件更新字段
func UpdateField(c *gin.Context) {
	req := make(map[string]interface{})

	req["uid"] = ""

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	dst := ".\\files\\" + req["uid"].(string) + ".json"

	fileBytes, err := ioutil.ReadFile(dst)

	if err != nil {
		log.Printf("unable to read the file. %v\n", err)
	}

	var fileMap map[string]interface{}

	err = json.Unmarshal(fileBytes, &fileMap)

	if err != nil {
		log.Printf("unmarshal the data failed. %v\n", err)
		return
	}
	var changed []string
	// 遍历寻找需要更新的字段并修改更新后的值，如果文件中没有则认为用户字段错误，跳过
	// 这里要求的是客户需要完整写出字段json格式的层级关系
	for tk, tv := range req {
		if tk == "uid" {
			// 不能修改uid字段
			continue
		}
		// 是否包含该一级字段
		if _, ok := fileMap[tk]; ok {
			fileMap[tk] = tv // 更新字段
			changed = append(changed, tk)
		}
	}

	jsonByte, _ := json.Marshal(fileMap)

	// 将更新后的信息重写并覆盖源文件
	err = ioutil.WriteFile(dst, jsonByte, 0644)

	if err != nil {
		log.Printf("write file failed. %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("%s have changed.", changed),
	})
}

// 个性化拉取
// 将用户的个性化参数持久化保存
func Pull(c *gin.Context) {

	var req Models.PerChoice

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	//for _, s := range req.Pull{
	//	fmt.Println(s)
	//	fmt.Println(s == "FirstName")
	//}

	// 读配置文件信息
	dst := ".\\files\\" + req.Uid + ".json"

	fileBytes, err := ioutil.ReadFile(dst)

	if err != nil {
		log.Printf("unable to read the file. %v\n", err)
	}

	var fileMap map[string]interface{}

	err = json.Unmarshal(fileBytes, &fileMap)

	if err != nil {
		log.Printf("unmarshal the data failed. %v\n", err)
		return
	}

	// 持久化个性参数数据
	db, err := bolt.Open(Utils.PULLFILE, 0644, nil)
	if err != nil {
		log.Printf("open or create  the db error. %v\n", err)
		return
	}

	defer db.Close() // 千万不能掉，否则连续请求就会失败

	// 保存从数据库中读取的信息
	var dataByte []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.PULLBUCKET))
		if b == nil {
			// 说明整个服务器是第一次被传个性化参数
			b, err = tx.CreateBucket([]byte(Utils.PULLBUCKET))
			if err != nil {
				log.Printf("Create the bucket failed. %v\n", err)
			}
		}

		if b != nil {
			dataByte = b.Get([]byte(req.Uid))
			pullfields := make(map[string]interface{})
			if dataByte == nil {
				// 说明之前没有该用户的个性化参数
				err = b.Put([]byte(req.Uid), Utils.Serialize(req.Pull))
				if err != nil {
					log.Printf("put the file into the db failed. %v\n", err)
				}

				for _, key := range req.Pull {
					fmt.Println(key)
					if _, ok := fileMap[key]; ok {
						fmt.Println("get it")
						pullfields[key] = fileMap[key]
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"code":   200,
					"info":   fmt.Sprintf("%v", pullfields),
					"config": fmt.Sprintf("%s", fileBytes),
				})

			} else {
				choice := Utils.Deserialize(dataByte)
				for _, key := range choice {
					if _, ok := fileMap[key]; ok {
						pullfields[key] = fileMap[key]
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"code":   200,
					"info":   fmt.Sprintf("%v", pullfields),
					"config": fmt.Sprintf("%s", fileBytes),
				})
			}
		}
		return nil
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
