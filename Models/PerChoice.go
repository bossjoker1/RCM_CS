package Models

// 用户自定义需求结构体

type PerChoice struct {
	// 用户id
	Uid int64 `json:"uid"`
	// 为了适应任意需求，这里采用的是动态类型
	Request map[string]interface{} `json:"request"`
}
