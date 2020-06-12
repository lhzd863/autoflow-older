package module

type MetaSystemUserBean struct {
	Id           string `json:"id"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	CreateTime   string `json:"createtime"`
	UpdateTime   string `json:"updatetime"`
	Enable       string `json:"enable"`
}

type MetaParaSystemUserAddBean struct {
	Id           string `json:"id"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	Enable       string `json:"enable"`
}

type MetaParaSystemUserUpdateBean struct {
	Id           string `json:"id"`
	Avatar       string `json:"avatar"`
	Introduction string `json:"introduction"`
	Enable       string `json:"enable"`
}

type MetaParaSystemUserGetBean struct {
	Id       string `json:"id"`
        UserName     string `json:"username"`
}

type MetaParaSystemUserRemoveBean struct {
	Id       string `json:"id"`
}
