package module

import (
	"os/exec"
)

type MetaJobMemBean struct {
	Id         string    `json:"id"`
	FlowId     string    `json:"flowid"`
	Sys        string    `json:"sys"`
	Job        string    `json:"job"`
	CreateTime string    `json:"cts"`
	UpdateTime string    `json:"uts"`
	Step       string    `json:"step"`
	Cmd        *exec.Cmd `josn:"cmd"`
	Enable     string    `json:"enable"`
}
