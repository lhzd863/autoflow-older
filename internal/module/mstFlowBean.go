package module

type MetaMstFlowBean struct {
	FlowId        string                 `json:"flowid"`
        MstId         string                 `json:"mstid"`
        ProcessNum    string                 `json:"processnum"`
	Ip            string                 `json:"ip"`
	Port          string                 `json:"port"`
        UpdateTime    string                 `json:"updatetime"`
}
