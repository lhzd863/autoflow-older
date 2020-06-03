package module

type MetaParaSystemWorkerRoutineJobRunningHeartAddBean struct {
        Id         string        `json:"id"`
	WorkerId    string        `json:"workerid"`
        Sys        string        `json:"sys"`
        Job        string        `json:"job"`
	Ip         string        `json:"ip"`
	Port       string        `json:"port"`
	UpdateTime string        `json:"updatetime"`
        StartTime  string        `json:"starttime"`
        Duration   string        `json:"duration"`
}
