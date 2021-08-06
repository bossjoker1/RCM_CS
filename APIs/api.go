package APIs

import (
	_ "RCM_CS/Utils"
	"encoding/json"
	"fmt"
	_ "github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func UploadByProperties(c *gin.Context) {
	// 获取客户端IP
	// clientIP := c.ClientIP()
	// 获取客户端port
	file, err := c.FormFile("file")
	// 获取用户传的uid 用于判断
	uid := c.PostForm("uid")
	fmt.Println("uid: ", uid)
	if err != nil {
		log.Printf("upload file failed. %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("ERROR: upload file failed. %s", err),
		})
		return
	}

	// 每个用户对应的专属目录
	_ = os.Mkdir(".\\files\\"+uid, 0666)

	var dst string

	if uid != "" {
		// 说明用户上传的是个性化配置文件
		// 命名为uid.properties 即与目录同名
		// 完整文件配置因为不知道用户上传的文件名是啥，所以无法限制

		// 给对应用户创建个人配置文件的目录
		dst = fmt.Sprintf(".\\files\\" + uid + "\\" + uid + ".properties")
		// 保存properties文件
		err = c.SaveUploadedFile(file, dst)
	} else {
		// 给对应用户创建个人配置文件的目录
		dst := fmt.Sprintf(".\\files\\default" + "\\" + file.Filename) // 按用户上传的文件保存
		// 保存properties文件
		err = c.SaveUploadedFile(file, dst)
	}

	if err != nil {
		log.Printf("save the file failed. %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("ERROR: save the  file failed. %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "upload the file successfully",
		"filepath": dst,
	})
}

// 用户指定字段更新或添加
func PersonalizedUpdate(c *gin.Context) {

	// 获得uid
	uid := c.Query("uid")

	// 按照properties文件格式，将kv全视为value型
	var req map[string]interface{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("bind the request json failed. %v\n", err)
	}

	// 读取文件配置
	// 这里认为每个目录对应一个配置文件
	// 所以就取第一个文件为配置文件读取

	viper.SetConfigName("config")
	viper.SetConfigType("properties")
	viper.AddConfigPath(".\\files\\" + uid)
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

	// 同时需要更新个性化参数配置文件的相关字段
	if _, err := os.Stat(".\\files\\" + uid + "\\" + uid + ".properties"); !os.IsNotExist(err) {
		// 读取总配置文件
		viper.SetConfigName(uid)
		viper.SetConfigType("properties")
		viper.AddConfigPath(".\\files\\" + uid)
		err = viper.ReadInConfig()
		if err != nil {
			log.Printf("read the properties file failed. %v\n", err)
			return
		}

		// 遍历所有的Key值，如果有则覆盖完成值的更新，如果没有添加
		for tk, tv := range req {
			if viper.Get(tk) != nil {
				viper.Set(tk, tv) //Set函数直接实现
			}
		}
		// 重新写入当前文件
		err = viper.WriteConfig() // 这个写入的时候会把所有大写转为小写再覆盖或者添加，不过不影响功能
		if err != nil {
			log.Printf("update the config file failed. %v\n", err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully update the config file.",
	})
}

//
//func PersonalizedPull2(c *gin.Context)  {
//
//	var req Models.PersonalizedFields
//
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		log.Printf("json bind failed. %v\n", err)
//		return
//	}
//
//	// 如果uid为空，则返回完整配置文件
//	if req.Uid == ""{
//		_, err := os.Open(".\\files\\"  + "\\" + filename)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError,
//				gin.H{
//					"msg": " file not existed.",
//				})
//			return
//		}
//
//		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
//		c.Writer.Header().Add("Content-Type", "application/octet-stream")
//
//		// 浏览器下载文件
//		c.File(".\\files\\" + clientIP + "\\" + filename)
//	}
//}

// 个性化拉取
func PersonalizedPull(c *gin.Context) {

	uid := c.Query("uid")

	var req []string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("json bind failed. %v\n", err)
		return
	}

	// 读取文件配置
	// 这里认为每个目录对应一个配置文件
	// 所以就取第一个文件为配置文件读取
	// 文件不存在说明用户还未上传过个性化配置
	// 用户的个性化文件格式为 uid.properties

	// 如果文件存在，则直接返回该文件的信息
	if _, err := os.Stat(".\\files\\" + uid + "\\" + uid + ".properties"); !os.IsNotExist(err) && len(req) == 0 {
		// 读取总配置文件
		viper.SetConfigName(uid)
		viper.SetConfigType("properties")
		viper.AddConfigPath(".\\files\\" + uid)
		err = viper.ReadInConfig()
		if err != nil {
			log.Printf("read the properties file failed. %v\n", err)
			return
		}
		jsonByte, _ := json.Marshal(viper.AllSettings())
		c.JSON(http.StatusOK, gin.H{
			"code":                200,
			"personalized_fields": fmt.Sprintf("%s", jsonByte),
		})
	} else {
		// 读取总配置文件
		viper.SetConfigName("config")
		viper.SetConfigType("properties")
		viper.AddConfigPath(".\\files\\" + uid + "\\")
		err = viper.ReadInConfig()
		if err != nil {
			log.Printf("read the properties file failed. %v\n", err)
			return
		}
		pullfields := make(map[string]interface{})
		for _, tk := range req {
			pullfields[tk] = viper.Get(tk)
		}

		// 创建个性化配置文件
		f, err := os.OpenFile(".\\files\\"+uid+"\\"+uid+".properties", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			log.Printf("create the personalized config file failed. %v\n", err)
		}
		// 新的文件流
		viper.SetConfigName(uid)
		viper.SetConfigType("properties")
		viper.AddConfigPath(".\\files\\" + uid)
		err = viper.ReadInConfig()
		if err != nil {
			log.Printf("read the properties file failed. %v\n", err)
			return
		}
		for tk, tv := range pullfields {
			viper.Set(tk, tv)
		}
		err = viper.WriteConfig()
		if err != nil {
			log.Printf("viper write the file failed. %v\n", err)
		}
		// 返回个性化参数
		jsonByte, _ := json.Marshal(viper.AllSettings())
		c.JSON(http.StatusOK, gin.H{
			"code":                200,
			"personalized_fields": fmt.Sprintf("%s", jsonByte),
		})
	}
}

func DownloadHandler(c *gin.Context) {
	var dst string

	// 改用query传多个参数
	filename := c.Query("filename")
	uid := c.Query("uid")

	if uid == "" {
		dst = ".\\files\\" + "default" + "\\" + filename
	} else {
		if filename == "" {
			// 如果没填文件名，则将uid指定目录下第一个文件返回
			files, _ := ioutil.ReadDir(".\\files\\" + uid)
			dst = ".\\files\\" + uid + "\\" + files[0].Name()
		} else {
			dst = ".\\files\\" + uid + "\\" + filename
		}
	}

	_, err := os.Open(dst)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"msg": " file not existed.",
			})
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	// 浏览器下载文件
	c.File(dst)

	return

}
