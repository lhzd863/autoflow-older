package apiserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"
)

type ResponseResourceWorker struct {
	sync.Mutex
}

func NewResponseResourceWorker() *ResponseResourceWorker {
	return &ResponseResourceWorker{}
}

func (rrs *ResponseResourceWorker) WorkerHeartAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaWorkerHeartAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()
	rrs.Lock()
	m := new(module.MetaWorkerMgrBean)
	v := bt.Get(p.Id)
	if v != nil {
		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
	} else {
		m.Id = p.Id
		m.WorkerId = p.WorkerId
		m.Ip = p.Ip
                m.Port = p.Port
		m.CurrentExecCnt = "0"
		m.CurrentSubmitCnt = "0"
		m.StartTime = p.StartTime
	}
	m.MaxCnt = p.MaxCnt
	m.Duration = p.Duration
	m.RunningCnt = p.RunningCnt

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	rrs.Unlock()
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerHeartRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaWorkerHeartGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	rrs.Lock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	err = bt.Remove(p.Id)
	rrs.Unlock()
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerHeartListHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaWorkerMgrBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v:%v.", m.WorkerId, m.Ip, m.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerHeartGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaWorkerHeartGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	rrs.Lock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	m := bt.Get(p.Id)
	if m != nil {
		v := new(module.MetaWorkerMgrBean)
		err := json.Unmarshal([]byte(m.(string)), &v)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		ise, _ := util.IsExpired(v.UpdateTime, timeStr, 300)
		if ise {
			glog.Glog(LogF, fmt.Sprintf("%v timeout %v:%v.", v.WorkerId, v.Ip, v.Port))
			bt.Remove(p.Id)
		} else {
			retlst = append(retlst, v)
		}
	}
	rrs.Unlock()
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerCntAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaWorkerMgrBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	p.UpdateTime = timeStr

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()
	fb0 := bt.Get(p.Id)
	fb := new(module.MetaWorkerMgrBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.CurrentExecCnt = p.CurrentExecCnt
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerMgrExecHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	rrs.Unlock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	strlist := bt.Scan()

	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaWorkerMgrBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			loc, _ := time.LoadLocation("Local")
			timeLayout := "2006-01-02 15:04:05"
			stheTime, _ := time.ParseInLocation(timeLayout, m.UpdateTime, loc)
			sst := stheTime.Unix()
			etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
			est := etheTime.Unix()
			if est-sst > 300 {
				glog.Glog(LogF, fmt.Sprintf("%v, %v:%v heart timeout.", m.WorkerId, m.Ip, m.Port))
				_ = bt.Remove(k1)
				continue
			}
			maxcnt, err := strconv.Atoi(m.MaxCnt)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("conv maxcnt fail.%v", err))
				continue
			}
			runningcnt, err := strconv.Atoi(m.RunningCnt)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("conv runningcnt fail.%v", err))
				continue
			}
			currentexeccnt, err := strconv.Atoi(m.CurrentExecCnt)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("conv currentcnt fail.%v", err))
				continue
			}
			currentsubmitcnt, err := strconv.Atoi(m.CurrentSubmitCnt)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("conv currentsubmitcnt fail.%v", err))
				continue
			}
			if 5*maxcnt <= runningcnt+currentexeccnt+currentsubmitcnt {
				glog.Glog(LogF, fmt.Sprintf("5*maxcnt(%v)<=runningcnt(%v)+currentcnt(%v)+currentsubmitcnt(%v).", maxcnt, runningcnt, currentexeccnt, currentsubmitcnt))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerRoutineJobRunningHeartListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_SLAVE_JOB_RUNNING_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemWorkerRoutineJobRunningHeartBean)
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

func (rrs *ResponseResourceWorker) WorkerRoutineJobRunningHeartGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemWorkerRoutineJobRunningHeartGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_SLAVE_JOB_RUNNING_HEART)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Id)
	if ib != nil {
		m := new(module.MetaSystemMstFlowRoutineJobRunningHeartBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceWorker) WorkerRoutineJobRunningHeartRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemWorkerRoutineJobRunningHeartRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_SLAVE_JOB_RUNNING_HEART)
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

func (rrs *ResponseResourceWorker) WorkerRoutineJobRunningHeartAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemWorkerRoutineJobRunningHeartAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_SLAVE_JOB_RUNNING_HEART)
	defer bt.Close()

	m := new(module.MetaSystemWorkerRoutineJobRunningHeartBean)
	m.Id = p.Id
	m.WorkerId = p.WorkerId
	m.Sys = p.Sys
	m.Job = p.Job
	m.StartTime = p.StartTime
	m.Ip = p.Ip
	m.Port = p.Port
	m.Duration = p.Duration
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr

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
