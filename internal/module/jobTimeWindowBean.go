package module

type MetaJobTimeWindowBean struct {
	Sys         string `json:"sys"`
	Job         string `json:"job"`
	Allow       string `json:"allow"`
	StartHour   int8   `json:"starthour"`
	EndHour     int8   `json:"endhour"`
	Description string `json:"description"`
	Enable      string `json:"enable"`
}
