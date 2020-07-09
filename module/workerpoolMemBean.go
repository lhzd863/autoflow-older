package module

import (
	"github.com/lhzd863/autoflow/workpool"
)

type MetaWorkerPoolMemBean struct {
	Id         string             `json:"id"`
	WorkerPool *workpool.WorkPool `json:"workerpool"`
	CreateTime string             `json:"cts"`
	Enable     string             `json:"enable"`
}
