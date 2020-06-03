package module

type MetaSystemImageBean struct {
	ImageId     string `json:"imageid"`
	Tag         string `json:"tag"`
	CreateTime  string `json:"createtime"`
	DbStore     string `json:"dbstore"`
	User        string `json:"user"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
