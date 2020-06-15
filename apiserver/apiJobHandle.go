package apiserver

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"
)

type ResponseResourceJob struct {
	sync.Mutex
}

func NewResponseResourceJob() *ResponseResourceJob {
	return &ResponseResourceJob{}
}

func (rrs *ResponseResourceJob) FlowJobListHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobGetHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	m := bt.Get(p.Sys + "." + p.Job)
	if m != nil {
		r := new(module.MetaJobBean)
		err := json.Unmarshal([]byte(m.(string)), &r)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, r)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobAddHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 != nil {
		glog.Glog(LogF, fmt.Sprintf("%v %v has exists.", p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v %v has exists.", p.Sys, p.Job), nil)
		return
	}
	mj := new(module.MetaJobBean)
	mj.Sys = p.Sys
	mj.Job = p.Job
	mj.Enable = p.Enable
	mj.TimeWindow = p.TimeWindow
	mj.RetryCnt = p.RetryCnt
	mj.Alert = p.Alert
	mj.TimeTrigger = p.TimeTrigger
	mj.JobType = p.JobType
	mj.Frequency = p.Frequency
	mj.CheckBatStatus = p.CheckBatStatus
	mj.Priority = p.Priority
	jsonstr, _ := json.Marshal(mj)
	err = bt.Set(mj.Sys+"."+mj.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobUpdateHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 != nil {
		mj := new(module.MetaJobBean)
		err := json.Unmarshal([]byte(fb0.(string)), &mj)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		mj.Enable = p.Enable
		mj.TimeWindow = p.TimeWindow
		mj.RetryCnt = p.RetryCnt
		mj.Alert = p.Alert
		mj.TimeTrigger = p.TimeTrigger
		mj.JobType = p.JobType
		mj.Frequency = p.Frequency
		mj.CheckBatStatus = p.CheckBatStatus
		mj.Priority = p.Priority
		mj.Status = p.Status
		jsonstr, _ := json.Marshal(mj)
		err = bt.Set(mj.Sys+"."+mj.Job, string(jsonstr))
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
			return
		}
	} else {
		glog.Glog(LogF, fmt.Sprintf("flow %v job %v %v no store data.", p.FlowId, p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flow %v job %v %v no store data.", p.FlowId, p.Sys, p.Job), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobRemoveHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	bt.Remove(p.Sys + "," + p.Job)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusGetPendingHandle(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusGetPendingBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	c, bid := rrs.CurrentStatusPendingOffset(len(strlist))
	retlst := make([]interface{}, 0)
	if bid == -1 {
                glog.Glog(LogF, fmt.Sprint("all ringid is working."))
		util.ApiResponse(response.ResponseWriter, 200, fmt.Sprintf("all ringid is working."), retlst)
		return
	} else {
		r := new(module.MetaRingPendingOffsetBean)
		r.Id = jobpara.Id
		r.RingId = fmt.Sprint(bid)
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		r.CreateTime = timeStr
		ringPendingSpool.Add(jobpara.Id, r)
	}
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobbn := new(module.MetaJobBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobbn)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobbn.Status != jobpara.Status && jobpara.Status != "ALL" {
				continue
			}
			if (c && jobpara.IsHash == "0") || jobpara.IsHash == "1" {
				jobnodeid := statusPendingHashRing.Get(jobbn.Job).Id
				glog.Glog(LogF, fmt.Sprintf("job hash info:%v %v", jobbn.Job, jobnodeid))
				if jobnodeid != bid {
					glog.Glog(LogF, fmt.Sprintf("local node id %v ,job %v  mapping node id %v is not local node.", bid, jobbn.Job, jobnodeid))
					continue
				}
			}
			retlst = append(retlst, jobbn)
		}
	}
	if len(retlst) == 0 {
		ringPendingSpool.Remove(jobpara.Id)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusGetGoHandle(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusGetGoBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	c, bid := rrs.CurrentStatusGoOffset(len(strlist))
	retlst := make([]interface{}, 0)
	if bid == -1 {
		util.ApiResponse(response.ResponseWriter, 200, "all ringid is working.", retlst)
		return
	} else {
		r := new(module.MetaRingPendingOffsetBean)
		r.Id = jobpara.Id
		r.RingId = fmt.Sprint(bid)
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		r.CreateTime = timeStr
		ringGoSpool.Add(jobpara.Id, r)
	}
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobbn := new(module.MetaJobBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobbn)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobbn.Status != jobpara.Status && jobpara.Status != "ALL" {
				continue
			}
			if (c && jobpara.IsHash == "0") || jobpara.IsHash == "1" {
				jobnodeid := statusGoHashRing.Get(jobbn.Job).Id
				glog.Glog(LogF, fmt.Sprintf("job hash info:%v %v", jobbn.Job, jobnodeid))
				if jobnodeid != bid {
					glog.Glog(LogF, fmt.Sprintf("local node id %v ,job %v  mapping node id %v is not local node.", bid, jobbn.Job, jobnodeid))
					continue
				}
			}
			retlst = append(retlst, jobbn)
		}
	}
	if len(retlst) == 0 {
		ringGoSpool.Remove(jobpara.Id)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobDependencyBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobdb := new(module.MetaJobDependencyBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobdb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobdb.Enable != "1" {
				continue
			}
			if jobdb.Sys == jobpara.Sys && jobdb.Job == jobpara.Job {
				bt1 := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
				defer bt1.Close()
				jobn := bt1.Get(jobdb.DependencySys + "." + jobdb.DependencyJob)
				bt1.Close()
				if jobn != nil {
					jobbn := new(module.MetaJobBean)
					err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
					if err != nil {
						glog.Glog(LogF, fmt.Sprint(err))
					}
					if jobbn.Status == util.STATUS_AUTO_SUCC || jobbn.Enable != "1" {
						continue
					}
					retlst = append(retlst, jobdb)
				}
			}
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobDependencyBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	if ib != nil {
		m := new(module.MetaJobDependencyBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()
	fb0 := bt.Get(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	fb := new(module.MetaJobDependencyBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.Sys+"."+p.Job+"."+p.DependencySys+"."+p.DependencyJob, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	err = bt.Remove(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobDependencyAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	jsonstr, _ := json.Marshal(p)
	err = bt.Set(p.Sys+"."+p.Job+"."+p.DependencySys+"."+p.DependencyJob, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamJobHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamJobBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 == nil {
		glog.Glog(LogF, fmt.Sprintf("%v.%v not exists.", p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v.%v not exists.", p.Sys, p.Job), nil)
		return
	}
	fb := new(module.MetaJobBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	if fb.Status != util.STATUS_AUTO_SUCC && fb.Status != util.STATUS_AUTO_READY {
		glog.Glog(LogF, fmt.Sprintf("%v.%v status %v not equal %v or %v .", p.Sys, p.Job, fb.Status, util.STATUS_AUTO_SUCC, util.STATUS_AUTO_READY))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v.%v status %v not equal %v or %v .", p.Sys, p.Job, fb.Status, util.STATUS_AUTO_SUCC, util.STATUS_AUTO_READY), nil)
		return
	}

	fb.Status = util.STATUS_AUTO_PENDING
	fb.EndTime = ""
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.StartTime = timeStr
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Sys+"."+fb.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamJobGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamJobListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobStreamBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if m.StreamSys == p.Sys && m.StreamJob == p.Job {
				retlst = append(retlst, m)
			}
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobStreamBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	if ib != nil {
		m := new(module.MetaJobStreamBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()
	fb0 := bt.Get(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	fb := new(module.MetaJobStreamBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.StreamSys+"."+p.StreamJob+"."+p.Sys+"."+p.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	err = bt.Remove(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStreamAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	m := new(module.MetaJobStreamBean)
	m.StreamSys = p.StreamSys
	m.StreamJob = p.StreamJob
	m.Sys = p.Sys
	m.Job = p.Job
	m.Description = p.Description
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.StreamSys+"."+m.StreamJob+"."+m.Sys+"."+m.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobTimeWindowBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job)
	if ib != nil {
		m := new(module.MetaJobTimeWindowBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()
	fb0 := bt.Get(p.Sys + "." + p.Job)
	fb := new(module.MetaJobTimeWindowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Allow = p.Allow
	fb.StartHour = p.StartHour
	fb.EndHour = p.EndHour
	fb.Enable = p.Enable
	fb.Description = p.Description

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.Sys+"."+p.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	err = bt.Remove(p.Sys + "." + p.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	m := new(module.MetaJobTimeWindowBean)
	m.Sys = p.Sys
	m.Job = p.Job
	m.Allow = p.Allow
	m.StartHour = p.StartHour
	m.EndHour = p.EndHour
	m.Description = p.Description
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Sys+"."+m.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobTimeWindowHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	jobtw := bt.Get(p.Sys + "." + p.Job)
	if jobtw != nil {
		jtw := new(module.MetaJobTimeWindowBean)
		err := json.Unmarshal([]byte(jobtw.(string)), &jtw)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		if jtw.Enable != "1" {
			util.ApiResponse(response.ResponseWriter, 200, "", retlst)
			return
		}
		bhour := jtw.StartHour
		ehour := jtw.EndHour
		timeStrHour := int8(time.Now().Hour())
		if timeStrHour >= bhour || timeStrHour <= ehour {
			util.ApiResponse(response.ResponseWriter, 200, "", retlst)
			return
		} else {
			retlst = append(retlst, jtw)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdGetAllHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdGetAllBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mjc := new(module.MetaJobCmdBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjc)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if mjc.Sys != jobpara.Sys || mjc.Job != jobpara.Job {
				continue
			}
			retlst = append(retlst, mjc)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job + "." + p.Step)
	if ib != nil {
		m := new(module.MetaJobCmdBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdRemoveHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdRemoveBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	bt.Remove(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Step)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobCmdBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	cb := new(module.MetaJobCmdBean)
	cb.Sys = p.Sys
	cb.Job = p.Job
	cb.Cmd = p.Cmd
	cb.Step = p.Step
	cb.Enable = p.Enable
	cb.Description = p.Description
	jsonstr, _ := json.Marshal(cb)

	err = bt.Set(p.Sys+"."+p.Job+"."+p.Step, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobCmdUpdateHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdUpdateBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	fb0 := bt.Get(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Step)
	fb := new(module.MetaJobCmdBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Type = jobpara.Type
	fb.Cmd = jobpara.Cmd
	fb.Step = jobpara.Step
	fb.Enable = jobpara.Enable
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Sys+"."+fb.Job+"."+fb.Step, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusUpdateSubmitHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdateSubmitBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusUpdateGoHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdateGoBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jobbn.SServer = jobpara.SServer
	jobbn.Sip = jobpara.Ip
	jobbn.Sport = jobpara.Port
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusUpdatePendingHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdatePendingBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jobbn.Sys = jobpara.Sys
	jobbn.Job = jobpara.Job
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusUpdateEndHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusUpdateEndBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	jobbn.EndTime = timeStr
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	jobpool.Remove(jobpara.FlowId + "," + jobpara.Sys + "," + jobpara.Job)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobStatusUpdateStartHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusUpdateStartBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	jobbn.StartTime = timeStr
	jobbn.EndTime = ""
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	mjm := new(module.MetaJobMemBean)
	mjm.FlowId = jobpara.FlowId
	mjm.Sys = jobbn.Sys
	mjm.Job = jobbn.Job
	mjm.CreateTime = timeStr
	mjm.Enable = "1"
	jobpool.Add(mjm.FlowId+","+mjm.Sys+","+mjm.Job, mjm)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterRemoveHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterRemoveBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	bt.Remove(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Key)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterGetHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterGetBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Sys != jobpara.Sys || m.Job != jobpara.Job || m.Key != jobpara.Key {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterGetAllHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterGetAllBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	f, err := rrf.flowInfo(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow info %v error.%v", jobpara.FlowId, err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow info %v error.%v", jobpara.FlowId, err), nil)
		return
	}
	rrs.Lock()
	defer rrs.Unlock()
	retmap := make(map[string]interface{})
	//system
	bt0 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt0.Close()
	strlist0 := bt0.Scan()
	bt0.Close()
	for _, v := range strlist0 {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			glog.Glog(LogF, fmt.Sprint(m))
			retmap[m.Key] = *m
		}
	}
	//flow
	bt1 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt1.Close()
	strlist1 := bt1.Scan()
	bt1.Close()
	for _, v := range strlist1 {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			glog.Glog(LogF, fmt.Sprint(m))
			retmap[m.Key] = *m
		}
	}
	p := new(module.MetaJobParameterBean)
	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW_CTS
	p.Val = f.CreateTime
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW_CTS] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW_RCT
	p.Val = f.RunContext
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW_RCT] = *p

	//job
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get parameter failed.%v", err), nil)
				return
			}
			if m.Sys != jobpara.Sys || m.Job != jobpara.Job {
				continue
			}
			retmap[m.Key] = *m
		}
	}
	//force
	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_SYS
	p.Val = jobpara.Sys
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_SYS] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_JOB
	p.Val = jobpara.Job
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_JOB] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW
	p.Val = jobpara.FlowId
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW] = *p

	retlst := make([]interface{}, 0)
	for _, v := range retmap {
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterListHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterListBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobParameterAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	m := new(module.MetaJobParameterBean)
	m.Sys = p.Sys
	m.Job = p.Job
	m.Key = p.Key
	m.Val = p.Val
	m.Description = p.Description
	m.Enable = p.Enable

	m.Type = "J"
	jsonstr, _ := json.Marshal(m)

	err = bt.Set(p.Sys+"."+p.Job+"."+p.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobParameterUpdateHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	fb0 := bt.Get(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Key)
	m := new(module.MetaJobParameterBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	m.Val = jobpara.Val
	m.Description = jobpara.Description
	m.Enable = jobpara.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Sys+"."+m.Job+"."+m.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) CurrentStatusPendingOffset(size int) (bool, int) {
	rrs.Lock()
	defer rrs.Unlock()
	t := -1
	c := false
	v := statusPendingOffset
	for i := 0; i < len(statusPendingHashRing.Resources); i++ {
		if statusPendingOffset == len(statusPendingHashRing.Resources)-1 {
			statusPendingOffset = 0
		} else {
			statusPendingOffset += 1
		}
		if rrs.IsExistPendingRingId(fmt.Sprint(v)) {
			v = statusPendingOffset
		} else {
			t = v
			break
		}
	}
	if size > 300 {
		c = true
	}
	return c, t
}

func (rrs *ResponseResourceJob) IsExistPendingRingId(ringid string) bool {
	for k := range ringPendingSpool.MemMap {
		v := ringPendingSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		if v.RingId == ringid {
			loc, _ := time.LoadLocation("Local")
			timeLayout := "2006-01-02 15:04:05"
			stheTime, _ := time.ParseInLocation(timeLayout, v.CreateTime, loc)
			sst := stheTime.Unix()
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
			est := etheTime.Unix()
			if est-sst > 300 {
				ringPendingSpool.Remove(k)
				break
			}
			return true
		}
	}
	return false
}

func (rrs *ResponseResourceJob) CurrentStatusGoOffset(size int) (bool, int) {
	rrs.Lock()
	defer rrs.Unlock()
	t := -1
	c := false
	v := statusGoOffset
	for i := 0; i < len(statusGoHashRing.Resources); i++ {
		if statusGoOffset == len(statusGoHashRing.Resources)-1 {
			statusGoOffset = 0
		} else {
			statusGoOffset += 1
		}
		if rrs.IsExistGoRingId(fmt.Sprint(v)) {
			v = statusGoOffset
		} else {
			t = v
			break
		}
	}

	if size > 300 {
		c = true
	}
	return c, t
}

func (rrs *ResponseResourceJob) IsExistGoRingId(ringid string) bool {
	for k := range ringGoSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		if v.RingId == ringid {
			loc, _ := time.LoadLocation("Local")
			timeLayout := "2006-01-02 15:04:05"
			stheTime, _ := time.ParseInLocation(timeLayout, v.CreateTime, loc)
			sst := stheTime.Unix()
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
			est := etheTime.Unix()
			if est-sst > 60 {
				ringGoSpool.Remove(k)
				break
			}
			return true
		}
	}
	return false
}

func (rrs *ResponseResourceJob) FlowJobLogListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobLogBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if m.Sys != p.Sys || m.Job != p.Job {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobLogGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Id)
	if ib != nil {
		m := new(module.MetaJobLogBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobLogRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	err = bt.Remove(p.Id)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobLogAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	m := new(module.MetaJobLogBean)
	m.Id = p.Id
	m.Sys = p.Sys
	m.Job = p.Job
	m.StartTime = p.StartTime
	m.EndTime = p.EndTime
	m.SServer = p.SServer
	m.Sip = p.Sip
	m.Sport = p.Sport
	m.Step = p.Step
	m.Content = make([]string, 0)
	m.ExitCode = p.ExitCode
	m.Cmd = p.Cmd
	for _, v := range p.Content {
		m.Content = append(m.Content, v)
	}
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(p.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceJob) FlowJobLogAppendHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrf.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	b := bt.Get(p.Id)
	m := new(module.MetaJobLogBean)
	if b != nil {
		err := json.Unmarshal([]byte(b.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("get in db error.%v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get in db error.%v", err), nil)
			return
		}
	} else {
		m.Id = p.Id
		m.Sys = p.Sys
		m.Job = p.Job
		m.StartTime = p.StartTime
		m.SServer = p.SServer
		m.Sip = p.Sip
		m.Sport = p.Sport
		m.Step = p.Step
		m.Cmd = p.Cmd
		m.Content = make([]string, 0)
	}
	for _, v := range p.Content {
		m.Content = append(m.Content, v)
	}
	m.EndTime = p.EndTime
	m.ExitCode = p.ExitCode
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(p.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}
