package module

type MetaJobWorkerBean struct {
	Id          string        `json:"id"`
	FlowId      string        `json:"flowid"`
	Sys         string        `json:"sys"`
	Job         string        `json:"job"`
	Retry       int           `json:"retry"`
	Alert       string        `json:"alert"`
	Status      string        `json:"status"`
	StartTime   string        `json:"stt"`
	EndTime     string        `json:"edt"`
	RunningTime string        `json:"rt"`
	RunningCmd  []interface{} `json:"rc"`
	Parameter   []interface{} `json:"parameter"`
	LeaderIp    string        `json:"lip"`
	LeaderPort  string        `json:"lport"`
	WorkerIp    string        `json:"wip"`
	WorkerPort  string        `json:"wport"`
	Cmd         []interface{} `json:"cmd"`
}
