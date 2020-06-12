package module

type MetaSystemRolePathBean struct {
	Id         string `json:"id"`
	Role       string `json:"role"`
	Path       string `json:"path"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
	Enable     string `json:"enable"`
}

type MetaParaSystemRolePathAddBean struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Remark string `json:"remark"`
	Enable string `json:"enable"`
}

type MetaParaSystemRolePathUpdateBean struct {
	Id         string `json:"id"`
	Role       string `json:"role"`
	Path       string `json:"path"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
	Enable     string `json:"enable"`
}

type MetaParaSystemRolePathRemoveBean struct {
	Id string `json:"id"`
}

type MetaParaSystemRolePathGetBean struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	Path string `json:"path"`
}
