package module

type MetaSystemUserTokenBean struct {
	Token string `json:"token"`
}

type MetaTokenCreateBean struct {
	UserName string `json:"username"`
	Hour     string `json:"hour"`
}
