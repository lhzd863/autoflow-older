package module

type MetaParaKVBean struct {
	FlowId         string        `json:"flowid"`
	Sys            string        `json:"sys"`
	Job            string        `json:"job"`
	Status         string        `json:"status"`
        IsHash         string        `json:"ishash"`
        Key            string        `json:"key"`
        Val            string        `json:"val"`
}
