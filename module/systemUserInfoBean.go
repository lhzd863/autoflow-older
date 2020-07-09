package module

type MetaSystemUserInfoBean struct {
	Id           string   `json:"id"`
	UserName     string   `json:"username"`
	Avatar       string   `json:"avatar"`
	Introduction string   `json:"introduction"`
	CreateTime   string   `json:"createtime"`
	UpdateTime   string   `json:"updatetime"`
	Role         []string `json:"role"`
	Enable       string   `json:"enable"`
}

type MetaParaSystemUserInfoGetBean struct {
	Token string `json:"token"`
}
