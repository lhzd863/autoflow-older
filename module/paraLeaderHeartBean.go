package module

type MetaParaLeaderHeartAddBean struct {
	Id         string `json:"id"`
	LeaderId   string `json:"leaderid"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	StartTime  string `json:"starttime"`
	UpdateTime string `json:"updatetime"`
	FlowNum    string `json:"flownum"`
	Duration   string `json:"duration"`
}

type MetaParaLeaderHeartGetBean struct {
	Id string `json:"id"`
}

type MetaSystemLeaderFlowRoutineJobRunningHeartBean struct {
	Id         string `json:"id"`
	LeaderId   string `json:"leaderid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	WorkerId   string `json:"workerid"`
	Sys        string `json:"sys"`
	Job        string `json:"job"`
	Wip        string `json:"wip"`
	Wport      string `json:"wport"`
	Mip        string `json:"mip"`
	Mport      string `json:"mport"`
	StartTime  string `json:"starttime"`
	UpdateTime string `json:"updatetime"`
	Duration   string `json:"duration"`
}

type MetaParaSystemLeaderFlowRoutineJobRunningHeartGetBean struct {
	LeaderId  string `json:"leaderid"`
	FlowId    string `json:"flowid"`
	RoutineId string `json:"routineid"`
	WorkerId  string `json:"workerid"`
	Sys       string `json:"sys"`
	Job       string `json:"job"`
}

type MetaParaSystemLeaderFlowRoutineJobRunningHeartAddBean struct {
	Id         string `json:"id"`
	LeaderId   string `json:"leaderid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	WorkerId   string `json:"workerid"`
	Sys        string `json:"sys"`
	Job        string `json:"job"`
	Wip        string `json:"wip"`
	Wport      string `json:"wport"`
	Mip        string `json:"mip"`
	Mport      string `json:"mport"`
	StartTime  string `json:"starttime"`
	UpdateTime string `json:"updatetime"`
	Duration   string `json:"duration"`
}

type MetaParaSystemLeaderFlowRoutineJobRunningHeartRemoveBean struct {
	LeaderId  string `json:"leaderid"`
	FlowId    string `json:"flowid"`
	RoutineId string `json:"routineid"`
	WorkerId  string `json:"workerid"`
	Sys       string `json:"sys"`
	Job       string `json:"job"`
}
