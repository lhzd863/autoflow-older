package module

type MetaParaFlowJobCmdGetBean struct {
	FlowId string `json:"flowid"`
	Sys    string `json:"sys"`
	Job    string `json:"job"`
	Step   string `json:"step"`
}
