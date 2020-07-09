package module

type MetaParaFlowJobParameterAddBean struct {
	FlowId      string `json:"flowid"`
	Sys         string `json:"sys"`
	Job         string `json:"job"`
	Key         string `json:"key"`
	Val         string `json:"val"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
