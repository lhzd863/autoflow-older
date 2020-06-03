package module

type MetaParaImageUpdateBean struct {
	ImageId     string `json:"imageid"`
        Tag         string `json:"tag"`
        DbStore     string `json:"dbstore"`
        Description string `json:"description"`
        Enable      string `json:"enable"`
}
