package module

type MetaMstFlowRoutineHeartBean struct {
	Id         string        `json:"id"`
	MstId      string        `json:"mstid"`
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
