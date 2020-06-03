package module

type MetaParaWorkerHeartAddBean struct {
        Id          string        `json:"id"`
        WorkerId     string        `json:"workerid"`
        Ip          string        `json:"ip"`
        Port        string        `json:"port"`
        MaxCnt      string        `json:"maxcnt"`
        RunningCnt  string        `json:"runningcnt"`
        CurrentCnt  string        `json:"currentcnt"`
        StartTime   string        `json:"starttime"`
        UpdateTime  string        `json:"updatetime"`
        Duration    string        `json:"duration"`
}
