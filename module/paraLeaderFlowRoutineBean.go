package module

type MetaParaLeaderFlowRoutineRemoveBean struct {
	LeaderId  string `json:"leaderid"`
	FlowId    string `json:"flowid"`
	RoutineId string `json:"routineid"`
}

type MetaParaLeaderFlowRoutineHeartAddBean struct {
	Id         string        `json:"id"`
	LeaderId   string        `json:"leaderid"`
	FlowId     string        `json:"flowid"`
	RoutineId  string        `json:"routineid"`
	Ip         string        `json:"ip"`
	Port       string        `json:"port"`
	StartTime  string        `json:"starttime"`
	UpdateTime string        `json:"updatetime"`
	Duration   string        `json:"duration"`
	JobNum     string        `json:"jobnum"`
	Lst        []interface{} `json:"lst"`
}

type MetaParaLeaderFlowRoutineHeartGetBean struct {
	Id string `json:"id"`
}

type MetaParaLeaderFlowRoutineHeartRemoveBean struct {
	Id string `json:"id"`
}
