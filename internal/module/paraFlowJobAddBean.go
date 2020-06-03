package module

type MetaParaFlowJobAddBean struct {
        FlowId         string `json:"flowid"`
	Sys            string `json:"sys"`
	Job            string `json:"job"`
        RunContext     string `json:"runcontext"`
        Status         string `json:"status"`
        JobType        string `json:"jobtype"`
        Description    string `json:"description"`
	Enable         string `json:"enable"`
	TimeWindow     string `json:"timewindow"`
	RetryCnt       string `json:"retrycnt"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
}
