package module

type MetaJobPoolBean struct {
	FlowId        string `json:"flowid"`
	Sys           string `json:"sys"`
	Job           string `json:"job"`
	StartTime     string `json:"starttime"`
	Priority      string `json:"priority"`
	Server        string `json:"server"`
	DynamicServer string `json:"dynamicserver"`
	Enable        string `json:"enable"`
}
