package module

type MetaSystemStatJobStatusBean struct {
	FlowId  string `json:"flowid"`
	Ready   string `json:"ready"`
	Pending string `json:"pending"`
	Submit  string `json:"submit"`
	Go      string `json:"go"`
	Succ    string `json:"succ"`
	Fail    string `json:"fail"`
	Running string `json:"running"`
}
