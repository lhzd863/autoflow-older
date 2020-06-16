package module

type MetaJobBean struct {
	Sys            string `json:"sys"`
	Job            string `json:"job"`
	RunContext     string `json:"runcontext"`
	Status         string `json:"status"`
	StartTime      string `json:"starttime"`
	EndTime        string `json:"endtime"`
	JobType        string `json:"jobtype"`
	Description    string `json:"description"`
	Enable         string `json:"enable"`
	SServer        string `json:"sserver"`
	Sip            string `json:"sip"`
	Sport          string `json:"sport"`
	MServer        string `json:"mserver"`
	Mip            string `json:"mip"`
	Mport          string `json:"mport"`
	RoutineId      string `json:"routineid"`
	TimeWindow     string `json:"timewindow"`
	DynamicServer  string `json:"dynamicserver"` //N:don't dynamic allocate server,other: Dynamic allocate server
	RetryCnt       string `json:"retrycnt"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
	RunningCmd     string `json:"runningcmd"`
}

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
	DynamicServer  string `json:"dynamicserver"`
	RetryCnt       string `json:"retrycnt"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
}

type MetaParaFlowJobGetBean struct {
	FlowId string `json:"flowid"`
	Sys    string `json:"sys"`
	Job    string `json:"job"`
}

type MetaParaFlowJobListBean struct {
	FlowId string `json:"flowid"`
}

type MetaParaFlowJobRemoveBean struct {
	FlowId string `json:"flowid"`
	Sys    string `json:"sys"`
	Job    string `json:"job"`
}

type MetaParaFlowJobUpdateBean struct {
	FlowId         string `json:"flowid"`
	Sys            string `json:"sys"`
	Job            string `json:"job"`
	RunContext     string `json:"runcontext"`
	Status         string `json:"status"`
	JobType        string `json:"jobtype"`
	Description    string `json:"description"`
	Enable         string `json:"enable"`
	TimeWindow     string `json:"timewindow"`
	DynamicServer  string `json:"dynamicserver"`
	RetryCnt       string `json:"retrycnt"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
}
