package module

type MetaParaFlowJobDependencyGetBean struct {
        FlowId        string `json:"flowid"`
	Sys           string `json:"sys"`
	Job           string `json:"job"`
	DependencySys string `json:"dependencysys"`
	DependencyJob string `json:"dependencyjob"`
}
