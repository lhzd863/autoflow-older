package module

type MetaLeaderHeartBean struct {
	Id         string        `json:"id"`
	LeaderId   string        `json:"leaderid"`
	Ip         string        `json:"ip"`
	Port       string        `json:"port"`
	StartTime  string        `json:"starttime"`
	UpdateTime string        `json:"updatetime"`
	FlowNum    string        `json:"flownum"`
	Duration   string        `json:"duration"`
	ProcessNum string        `json:"processnum"`
	Lst        []interface{} `json:"lst"`
}

type MetaLeaderFlowRoutineHeartBean struct {
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

type MetaLeaderFlowBean struct {
	FlowId     string `json:"flowid"`
	LeaderId   string `json:"leaderid"`
	ProcessNum string `json:"processnum"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	UpdateTime string `json:"updatetime"`
}
