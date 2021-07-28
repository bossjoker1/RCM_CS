package Models

// 用户自定义需求结构体

type PerChoice struct {
	// 用户id
	Uid string `json:"uid"`
	// 字段名集合
	Pull []string `json:"pull"`
}
