package module

type MetaSystemUserRoleBean struct {
	Id         string `json:"id"`
	UserName   string `json:"username"`
	Role       string `json:"role"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
	Enable     string `json:"enable"`
}

type MetaParaSystemUserRoleAddBean struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
	Remark   string `json:"remark"`
	Enable   string `json:"enable"`
}

type MetaParaSystemUserRoleGetBean struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

type MetaParaSystemUserRoleRemoveBean struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

type MetaParaSystemUserRoleUpdateBean struct {
	Id         string `json:"id"`
	UserName   string `json:"username"`
	Role       string `json:"role"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
	Enable     string `json:"enable"`
}
