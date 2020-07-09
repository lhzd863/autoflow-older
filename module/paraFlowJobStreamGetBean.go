package module

type MetaParaFlowJobStreamGetBean struct {
	FlowId    string `json:"flowid"`
	StreamSys string `json:"streamsys"`
	StreamJob string `json:"streamjob"`
	Sys       string `json:"sys"`
	Job       string `json:"job"`
}
