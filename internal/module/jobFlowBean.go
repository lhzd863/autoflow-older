package module

type MetaJobFlowBean struct {
	FlowId     string                 `json:"flowid"`
	ImageId    string                 `json:"imageid"`
	RunContext string                 `json:"runcontext"`
	User       string                 `json:"user"`
	CreateTime string                 `json:"createtime"`
	StartTime  string                 `json:"starttime"`
	EndTime    string                 `json:"endtime"`
	Status     string                 `json:"status"`
	DbStore    string                 `json:"dbstore"`
	ParaMap    map[string]interface{} `json:"paramap"`
	Elapsed    string                 `json:"elapsed"`
	KeepPeriod string                 `json:"keepPeriod"`
	Ip         string                 `json:"ip"`
	Port       string                 `json:"port"`
	HomeDir    string                 `json:"homedir"`
	MstId      string                 `json:"mstid"`
	Enable     string                 `json:"enable"`
	ProcessNum string                 `json:"processnum"`
}
