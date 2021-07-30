package APIs

import (
	"RCM_CS/Models"
	"RCM_CS/Utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func UploadByProperties(c *gin.Context) {
	// 获取客户端IP
	//clientIP := c.ClientIP()
	// 获取客户端port
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("upload file failed. %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("ERROR: upload file failed. %s", err),
		})
		return
	}
	// 给对应用户创建个人配置文件的目录
	_ = os.Mkdir(".\\files\\files1", 0666)
	dst := fmt.Sprintf(".\\files\\files1\\" + file.Filename)
	// 保存properties文件
	err = c.SaveUploadedFile(file, dst)

	if err != nil {
		log.Printf("save the file failed. %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "upload the file successfully",
		"filepath": dst,
	})
}

// 用户指定字段更新或添加
func PersonalizedUpdate(c *gin.Context) {
	// 按照properties文件格式，将kv全视为value型
	var req map[string]interface{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind the request json failed. %v\n", err)
	}
	// 读取文件配置
	viper.SetConfigName("test")
	viper.SetConfigType("properties")
	viper.AddConfigPath(".\\files")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("read the properties file failed. %v\n", err)
		return
	}

	// 遍历所有的Key值，如果有则覆盖完成值的更新，如果没有添加
	for tk, tv := range req {
		fmt.Println("key: ", tk, "value: ", tv)
		viper.Set(tk, tv) //Set函数直接实现
	}
	// 重新写入当前文件
	err = viper.WriteConfig() // 这个写入的时候会把所有大写转为小写再覆盖或者添加，不过不影响功能
	if err != nil {
		log.Printf("update the config file failed. %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully update the config file.",
	})
}

// 个性化拉取
func PersonalizedPull(c *gin.Context) {

	var req Models.PersonalizedFields

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	// 读取文件配置
	viper.SetConfigName("test")
	viper.SetConfigType("properties")
	viper.AddConfigPath(".\\files")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("read the properties file failed. %v\n", err)
		return
	}

	jsonByte, err := json.Marshal(viper.AllSettings())
	if err != nil {
		log.Printf("marshal to json failed. %v\n", err)
	}

	// 持久化个性参数数据
	db, err := bolt.Open(Utils.PULLFILE, 0644, nil)
	if err != nil {
		log.Printf("open or create  the db error. %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "do not have the personalized pull_file.",
		})
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
			dataByte = b.Get([]byte("test"))
			pullfields := make(map[string]interface{})
			if dataByte == nil || req.Fields != nil {
				// 说明之前没有该用户的个性化参数
				err = b.Put([]byte("test"), Utils.Serialize(req.Fields))
				if err != nil {
					log.Printf("put the file into the db failed. %v\n", err)
				}

				for _, key := range req.Fields {
					fmt.Println(key)
					if v := viper.Get(key); v != nil {
						pullfields[key] = v
					}
				}

			} else {
				choice := Utils.Deserialize(dataByte)
				for _, key := range choice {
					if v := viper.Get(key); v != nil {
						pullfields[key] = v
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"code":   200,
				"info":   fmt.Sprintf("%v", pullfields),
				"config": fmt.Sprintf("%s", jsonByte),
			})
		}

		return nil
	})

}

func DownloadHandler(c *gin.Context) {
	filename := c.Param("filename")

	_, err := os.Open(".\\files\\" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"msg": " file not existed.",
			})
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	// 浏览器下载文件
	c.File(".\\files\\" + filename)

	return

}
