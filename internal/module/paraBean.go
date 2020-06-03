package module

type MetaParaBean struct {
	FlowId         string        `json:"flowid"`
	Sys            string        `json:"sys"`
	Job            string        `json:"job"`
	RunContext     string        `json:"runcontext"`
	Enable         string        `json:"enable"`
	SServer        string        `json:"sserver"`
	MServer        string        `json:"mstid"`
	Ip             string        `json:"ip"`
	Port           string        `json:"port"`
	Sip            string        `json:"sip"`
	Sport          string        `json:"sport"`
	Mip            string        `json:"mip"`
	Mport          string        `json:"mport"`
	TimeWindow     string        `json:"timewindow"`
	RetryCnt       string        `json:"retrycnt"`
	Alert          string        `json:"alert"`
	TimeTrigger    string        `json:"timetrigger"`
	JobType        string        `json:"jobtype"`
	Frequency      string        `json:"frequency"`
	Status         string        `json:"status"`
	StartTime      string        `json:"starttime"`
	EndTime        string        `json:"endtime"`
	RunTime        string        `json:"runtime"`
	CheckBatStatus string        `json:"checkbatstatus"`
	Priority       string        `json:"priority"`
	RunningCmd     string        `json:"runningcmd"`
	Cmd            string        `json:"cmd"`
	Type           string        `json:"type"`
	Step           string        `json:"step"`
	IsHash         string        `json:"ishash"`
	CodeType       string        `json:"codetype"`
	Id             string        `json:"id"`
	Key            string        `json:"key"`
	Val            string        `json:"val"`
	Description    string        `json:"description"`
	DbStore        string        `json:"dbstore"`
	HomeDir        string        `json:"homedir"`
	MstId          string        `json:"mstid"`
	ProcessNum     string        `json:"processnum"`
	RoutineId      string        `json:"routineid"`
	JobList        []interface{} `json:"joblist"`
}
