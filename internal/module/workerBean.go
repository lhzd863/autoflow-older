package module

type MetaWorkerHeartBean struct {
	Id             string `json:"id"`
	WorkerId       string `json:"workerid"`
	Ip             string `json:"ip"`
	Port           string `json:"port"`
	MaxCnt         string `json:"maxcnt"`
	RunningCnt     string `json:"runningcnt"`
	CurrentExecCnt string `json:"currentexeccnt"`
	StartTime      string `json:"starttime"`
	UpdateTime     string `json:"updatetime"`
	Duration       string `json:"duration"`
}

type MetaWorkerMgrBean struct {
	Id               string `json:"id"`
	WorkerId         string `json:"workerid"`
	Ip               string `json:"ip"`
	Port             string `json:"port"`
	MaxCnt           string `json:"maxcnt"`
	RunningCnt       string `json:"runningcnt"`
	CurrentExecCnt   string `json:"currentexeccnt"`
	CurrentSubmitCnt string `json:"currentsubmitcnt"`
	StartTime        string `json:"starttime"`
	UpdateTime       string `json:"updatetime"`
	Duration         string `json:"duration"`
}
