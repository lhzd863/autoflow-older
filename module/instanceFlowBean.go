package module

type MetaInstanceFlowBean struct {
	FlowId     string `json:"flowid"`
	LeaderId   string `json:"leaderid"`
	RoutineId  string `json:"routineid"`
	ProcessNum string `json:"processnum"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
}
