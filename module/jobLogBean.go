package module

type MetaJobLogBean struct {
	Id        string   `json:"id"`
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
