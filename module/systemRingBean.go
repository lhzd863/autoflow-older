package module

type MetaSystemRingBean struct {
	Id         string `json:"id"`
	LeaderId   string `json:"leaderid"`
	FlowId     string `json:"flowid"`
	RoutineId  string `json:"routineid"`
	CreateTime string `json:"createtime"`
	RingHash   string `json:"ringhash"`
}
