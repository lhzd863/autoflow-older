package module

type MetaSystemUserInfoBean struct {
	UserName     string   `json:"username"`
	Avatar       string   `json:"avatar"`
	Introduction string   `json:"introduction"`
	Role         []string `json:"role"`
}

type MetaParaSystemUserInfoGetBean struct {
	Token string `json:"token"`
}
