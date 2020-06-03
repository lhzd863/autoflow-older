package module

type MetaParaFlowJobStatusUpdateGoBean struct {
	FlowId         string        `json:"flowid"`
	Sys            string        `json:"sys"`
	Job            string        `json:"job"`
	Status         string        `json:"status"`
        SServer        string        `json:"sserver"`
        Ip             string        `json:"ip"`
        Port           string        `json:"port"`
}
