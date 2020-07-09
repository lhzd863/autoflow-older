package module

type MetaParaSystemMstFlowRoutineJobRunningHeartAddBean struct {
	Id         string `json:"id"`
	MstId      string `json:"mstid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	WorkerId   string `json:"workerid"`
	Sys        string `json:"sys"`
	Job        string `json:"job"`
	Sip        string `json:"sip"`
	Sport      string `json:"sport"`
	Mip        string `json:"mip"`
	Mport      string `json:"mport"`
	StartTime  string `json:"starttime"`
	UpdateTime string `json:"updatetime"`
	Duration   string `json:"duration"`
}
