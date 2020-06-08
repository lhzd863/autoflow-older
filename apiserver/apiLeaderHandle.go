package apiserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"
)

type ResponseResourceLeader struct {
}

func NewResponseResourceLeader() *ResponseResourceLeader {
	return &ResponseResourceLeader{}
}

func (rrs *ResponseResourceLeader) MstHeartAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstHeartAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	m := new(module.MetaMstHeartBean)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Id = p.Id
	m.MstId = p.MstId
	m.Ip = p.Ip
	m.Port = p.Port
	m.StartTime = p.StartTime
	m.FlowNum = p.FlowNum
	m.Duration = p.Duration

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstHeartListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaMstHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}

			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 600)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", m.MstId, m.Ip, m.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstHeartGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstHeartGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	m := bt.Get(p.Id)
	if m != nil {
		v := new(module.MetaMstHeartBean)
		err := json.Unmarshal([]byte(m.(string)), &v)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstHeartRemoveHandler(request *restful.Request, response *restful.Response) {
	m := new(module.MetaParaMstHeartRemoveBean)
	err := request.ReadEntity(&m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	err = bt.Remove(m.Id)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineHeartAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowRoutineHeartAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.MstId) == 0 || len(p.FlowId) == 0 || len(p.RoutineId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
	defer bt.Close()

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m := new(module.MetaMstFlowRoutineHeartBean)
	m.Id = p.Id
	m.MstId = p.MstId
	m.FlowId = p.FlowId
	m.RoutineId = p.RoutineId
	m.Ip = p.Ip
	m.Port = p.Port
	m.UpdateTime = timeStr
	m.StartTime = p.StartTime
	m.Lst = p.Lst
	m.JobNum = p.JobNum
	m.Duration = p.Duration

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineHeartListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaMstFlowRoutineHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}

			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v %v %v:%v.", m.MstId, m.FlowId, m.RoutineId, m.Ip, m.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineHeartGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowRoutineHeartGetBean)
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
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	v := bt.Get(p.Id)
	if v != nil {
		m := new(module.MetaMstFlowRoutineHeartBean)
		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineHeartRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowRoutineHeartRemoveBean)
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
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
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

func (rrs *ResponseResourceLeader) MstFlowRoutineJobRunningHeartListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOB_RUNNING_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemMstFlowRoutineJobRunningHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("MstId %v(%v:%v) FlowId %v RoutineId %v %v %v timeout on %v(%v:%v).", m.MstId, m.Mip, m.Mport, m.FlowId, m.RoutineId, m.Sys, m.Job, m.WorkerId, m.Sip, m.Sport))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineJobRunningHeartGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemMstFlowRoutineJobRunningHeartGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOB_RUNNING_HEART)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.MstId + "." + p.FlowId + "." + p.RoutineId + "." + p.WorkerId + "." + p.Sys + "." + p.Job)
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

func (rrs *ResponseResourceLeader) MstFlowRoutineJobRunningHeartRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemMstFlowRoutineJobRunningHeartRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOB_RUNNING_HEART)
	defer bt.Close()

	err = bt.Remove(p.MstId + "." + p.FlowId + "." + p.RoutineId + "." + p.WorkerId + "." + p.Sys + "." + p.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceLeader) MstFlowRoutineJobRunningHeartAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemMstFlowRoutineJobRunningHeartAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOB_RUNNING_HEART)
	defer bt.Close()

	m := new(module.MetaSystemMstFlowRoutineJobRunningHeartBean)
	m.Id = p.Id
	m.MstId = p.MstId
	m.FlowId = p.FlowId
	m.RoutineId = p.RoutineId
	m.WorkerId = p.WorkerId
	m.Sys = p.Sys
	m.Job = p.Job
	m.StartTime = p.StartTime
	m.Sip = p.Sip
	m.Sport = p.Sport
	m.Mip = p.Mip
	m.Mport = p.Mport
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
