package module

type MetaJobLogBean struct {
	Id        string   `json:"id"`
	Sys       string   `json:"sys"`
	Job       string   `json:"job"`
	Step      string   `json:"step"`
	WServer   string   `json:"wserver"`
	WIp       string   `json:"wip"`
	WPort     string   `json:"wport"`
	StartTime string   `json:"starttime"`
	EndTime   string   `json:"endtime"`
	Content   []string `json:"content"`
	ExitCode  string   `json:"exitcode"`
	Cmd       string   `json:"cmd"`
}
