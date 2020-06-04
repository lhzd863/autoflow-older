package module

type MetaInstanceFlowBean struct {
	FlowId     string `json:"flowid"`
	MstId      string `json:"mstid"`
	RoutineId  string `json:"routineid"`
	ProcessNum string `json:"processnum"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
}
