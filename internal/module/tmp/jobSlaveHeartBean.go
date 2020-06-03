package module

type MetaSlaveHeartBean struct {
	SlaveId    string        `json:"slaveid"`
	Ip         string        `json:"ip"`
	Port       string        `json:"port"`
	UpdateTime string        `json:"updatetime"`
	MaxCnt     string        `json:"maxcnt"`
	RunningCnt string        `json:"runningcnt"`
	CurrentCnt string        `json:"currentcnt"`
	JobList    []interface{} `json:"joblist"`
	Enable     string        `json:"enable"`
}
