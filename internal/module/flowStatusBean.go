package module

type MetaFlowStatusBean struct {
	FlowId         string `json:"flowid"`
	MstId          string `json:"mstid"`
	Ip             string `json:"ip"`
	Port           string `json:"port"`
	WorkPoolStatus string `json:"workpoolstatus"`
	MyWorkCnt      string `json:"myworkcnt"`
        Enable         string `json:"enable"`
}
