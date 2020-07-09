package module

type MetaParaFlowUpdateBean struct {
	FlowId     string `json:"flowid"`
	RunContext string `json:"runcontext"`
	DbStore    string `json:"dbstore"`
	HomeDir    string `json:"homedir"`
	Enable     string `json:"enable"`
}
