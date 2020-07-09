package module

type MetaParaMstHeartAddBean struct {
	Id         string `json:"id"`
	MstId      string `json:"mstid"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	StartTime  string `json:"starttime"`
	UpdateTime string `json:"updatetime"`
	FlowNum    string `json:"flownum"`
	Duration   string `json:"duration"`
}
