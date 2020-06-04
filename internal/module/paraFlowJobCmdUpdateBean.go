package module

type MetaParaFlowJobCmdUpdateBean struct {
	FlowId      string `json:"flowid"`
	Sys         string `json:"sys"`
	Job         string `json:"job"`
	Type        string `json:"type"`
	Step        string `json:"step"`
	Cmd         string `json:"cmd"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
