package module

type MetaJobStreamBean struct {
	StreamSys   string `json:"streamsys"`
	StreamJob   string `json:"streamjob"`
	Sys         string `json:"sys"`
	Job         string `json:"job"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
