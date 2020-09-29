package module

type MetaJobBean struct {
	Sys            string        `json:"sys"`
	Job            string        `json:"job"`
	RunContext     string        `json:"runcontext"`
	Status         string        `json:"status"`
	StartTime      string        `json:"starttime"`
	EndTime        string        `json:"endtime"`
	JobType        string        `json:"jobtype"`
	Description    string        `json:"description"`
	Enable         string        `json:"enable"`
	WServer        string        `json:"wserver"` //worker server
	Wip            string        `json:"wip"`
	Wport          string        `json:"wport"`
	Mserver        string        `json:"mserver"` //main server
	Mip            string        `json:"mip"`
	Mport          string        `json:"mport"`
	RoutineId      string        `json:"routineid"`
	TimeWindow     string        `json:"timewindow"`
	DynamicServer  string        `json:"dynamicserver"` //N:don't dynamic allocate server,other: Dynamic allocate server
	Retry          string        `json:"retry"`
	Alert          string        `json:"alert"`
	TimeTrigger    string        `json:"timetrigger"`
	Frequency      string        `json:"frequency"`
	CheckBatStatus string        `json:"checkbatstatus"`
	Priority       string        `json:"priority"`
	RunningCmd     string        `json:"runningcmd"`
	Server         []interface{} `json:"server"`
	Timeout        string        `json:"timeout"`
}

type MetaParaFlowJobAddBean struct {
	FlowId         string `json:"flowid"`
	Sys            string `json:"sys"`
	Job            string `json:"job"`
	WServer        string `json:"wserver"` //worker server
	Wip            string `json:"wip"`
	Wport          string `json:"wport"`
	RunContext     string `json:"runcontext"`
	Status         string `json:"status"`
	JobType        string `json:"jobtype"`
	Description    string `json:"description"`
	Enable         string `json:"enable"`
	TimeWindow     string `json:"timewindow"`
	DynamicServer  string `json:"dynamicserver"`
	Retry          string `json:"retry"`
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
	WServer        string `json:"wserver"`
	Wip            string `json:"wip"`
	Wport          string `json:"wport"`
	RunContext     string `json:"runcontext"`
	Status         string `json:"status"`
	JobType        string `json:"jobtype"`
	Description    string `json:"description"`
	Enable         string `json:"enable"`
	TimeWindow     string `json:"timewindow"`
	DynamicServer  string `json:"dynamicserver"`
	Retry          string `json:"retry"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
}

type MetaParaFlowJobStatus2ServerBean struct {
	FlowId string `json:"flowid"`
	Sys    string `json:"sys"`
	Job    string `json:"job"`
	Server string `json:"server"`
	Status string `json:"status"`
}

type MetaJobServerBean struct {
	Server string `json:"server"`
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	Type   string `json:"type"`
}
