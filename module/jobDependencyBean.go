package module

type MetaJobDependencyBean struct {
	Sys           string `json:"sys"`
	Job           string `json:"job"`
	DependencySys string `json:"dependencysys"`
	DependencyJob string `json:"dependencyjob"`
	Description   string `json:"description"`
	Enable        string `json:"enable"`
}
