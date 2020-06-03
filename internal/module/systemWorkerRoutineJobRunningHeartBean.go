package module

type MetaSystemWorkerRoutineJobRunningHeartBean struct {
        Id         string        `json:"id"`
	WorkerId    string        `json:"workerid"`
        Sys        string        `json:"sys"`
        Job        string        `json:"job"`
	Ip         string        `json:"ip"`
	Port       string        `json:"port"`
        StartTime  string        `json:"starttime"`
	UpdateTime string        `json:"updatetime"`
        Duration   string        `json:"duration"`
}
