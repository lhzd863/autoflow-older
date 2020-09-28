package module

type MetaParaFlowJobStatusUpdateGoBean struct {
	FlowId  string `json:"flowid"`
	Sys     string `json:"sys"`
	Job     string `json:"job"`
	Status  string `json:"status"`
	WServer string `json:"wserver"` //worker server
	WIp     string `json:"wip"`
	WPort   string `json:"wport"`
}
