package Models

// 用户自定义需求结构体

type PersonalizedFields struct {
	//用户id, 如果为空则说明请求默认配置
	Uid string `json:"uid"`
	// 个性化拉取的字段
	Fields []string `json:"fields"`
}

type Uid struct {
	Uid string `json:"uid"`
}
