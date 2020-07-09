package module

type RetBean struct {
	Status_Txt  string      `json:"status_txt"`
	Status_Code int         `json:"status_code"`
	Data        interface{} `json:"data"`
}
