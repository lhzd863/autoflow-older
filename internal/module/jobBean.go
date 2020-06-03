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
	RetryCnt       string `json:"retrycnt"`
	Alert          string `json:"alert"`
	TimeTrigger    string `json:"timetrigger"`
	Frequency      string `json:"frequency"`
	CheckBatStatus string `json:"checkbatstatus"`
	Priority       string `json:"priority"`
	RunningCmd     string `json:"runningcmd"`
}
