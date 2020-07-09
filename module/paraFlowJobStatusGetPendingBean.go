package module

type MetaParaFlowJobStatusGetPendingBean struct {
	Id     string `josn:"id"`
	FlowId string `json:"flowid"`
	Sys    string `json:"sys"`
	Job    string `json:"job"`
	Status string `json:"status"`
	IsHash string `json:"ishash"`
}
