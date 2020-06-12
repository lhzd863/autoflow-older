package module

type MetaSystemUserPasswordBean struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type MetaParaSystemUserPasswordChangeBean struct {
	UserName    string `json:"username"`
	OldPassword string `json:"oldpassword"`
	Password    string `json:"password"`
}
