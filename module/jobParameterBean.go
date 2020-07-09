package module

type MetaJobParameterBean struct {
	FlowId      string `json:"flowid"`
	Sys         string `json:"sys"`
	Job         string `json:"job"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	Val         string `json:"val"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
