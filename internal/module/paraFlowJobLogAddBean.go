package module

type MetaParaFlowJobLogAddBean struct {
	Id        string   `json:"id"`
	FlowId    string   `json:"flowid"`
	Sys       string   `json:"sys"`
	Job       string   `json:"job"`
	Step      string   `json:"step"`
	SServer   string   `json:"sserver"`
	Sip       string   `json:"sip"`
	Sport     string   `json:"sport"`
	StartTime string   `json:"starttime"`
	EndTime   string   `json:"endtime"`
	Content   []string `json:"content"`
	ExitCode  string   `json:"exitcode"`
	Cmd       string   `json:"cmd"`
}
