package module

type MetaParaFlowJobLogAddBean struct {
	Id        string   `json:"id"`
	FlowId    string   `json:"flowid"`
	Sys       string   `json:"sys"`
	Job       string   `json:"job"`
	Step      string   `json:"step"`
	WServer   string   `json:"wserver"` //worker server
	Wip       string   `json:"wip"`
	Wport     string   `json:"wport"`
	StartTime string   `json:"starttime"`
	EndTime   string   `json:"endtime"`
	Content   []string `json:"content"`
	ExitCode  string   `json:"exitcode"`
	Cmd       string   `json:"cmd"`
}
