package module

type MetaParaLeaderFlowStartBean struct {
	LeaderId   string `json:"leaderid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	ProcessNum string `json:"processnum"`
}

type MetaParaLeaderFlowStopBean struct {
	LeaderId   string `json:"leaderid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	ProcessNum string `json:"processnum"`
}

type MetaParaLeaderHeartRemoveBean struct {
	Id string `json:"id"`
}

type MetaParaLeaderFlowRemoveBean struct {
	FlowId   string `json:"flowid"`
	LeaderId string `json:"leaderid"`
}
