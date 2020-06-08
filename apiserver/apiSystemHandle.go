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

type ResponseResourceSystem struct {
}

func NewResponseResourceSystem() *ResponseResourceSystem {
	return &ResponseResourceSystem{}
}

func (rrs *ResponseResourceSystem) SystemParameterListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaParaFlowBean)
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

func (rrs *ResponseResourceSystem) SystemParameterGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemParameterGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Key)
	if ib != nil {
		m := new(module.MetaParaFlowBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemParameterUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemParameterUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt.Close()
	fb0 := bt.Get(p.Key)
	fb := new(module.MetaParaFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Val = p.Val
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemParameterRemoveHandler(request *restful.Request, response *restful.Response) {
	m := new(module.MetaParaSystemParameterRemoveBean)
	err := request.ReadEntity(&m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt.Close()

	err = bt.Remove(m.Key)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemParameterAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemParameterAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	m := new(module.MetaJobParameterBean)
	m.Type = "S"
	m.Key = p.Key
	m.Val = p.Val
	m.Description = p.Description
	m.Enable = p.Enable
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt.Close()

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(p.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SysListPortHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
	defer bt.Close()

	strlist := bt.Scan()
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			mmf := new(module.MetaMstFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmf)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			ism, _ := rrf.IsExpiredMst(mmf.MstId)
			if ism {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", mmf.MstId, mmf.Ip, mmf.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, mmf)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) JobPoolAddHandler(request *restful.Request, response *restful.Response) {
	jpbean := new(module.MetaJobPoolBean)
	err := request.ReadEntity(&jpbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jpbean.FlowId) == 0 || len(jpbean.Sys) == 0 || len(jpbean.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("flowid sys or job missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flowid sys job missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOBSPOOL)
	defer bt.Close()

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	jpbean.StartTime = timeStr
	jpbean.Enable = "1"

	jsonstr, _ := json.Marshal(jpbean)
	err = bt.Set(jpbean.FlowId+"."+jpbean.Sys+"."+jpbean.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) JobPoolGetHandler(request *restful.Request, response *restful.Response) {
	jpbean := new(module.MetaJobPoolBean)
	err := request.ReadEntity(&jpbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(jpbean.FlowId) == 0 || len(jpbean.Sys) == 0 || len(jpbean.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("flowid sys or job missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flowid sys job missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOBSPOOL)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	jp := bt.Get(jpbean.FlowId + "." + jpbean.Sys + "." + jpbean.Job)
	if jp != nil {
		mjpb := new(module.MetaJobPoolBean)
		err := json.Unmarshal([]byte(jp.(string)), &mjpb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, mjpb)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) JobPoolListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOBSPOOL)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			mjpb := new(module.MetaJobPoolBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjpb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(mjpb.StartTime, timeStr, 600)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v %v %v exist job pool timeout,will remove from pool.", mjpb.FlowId, mjpb.Sys, mjpb.Job))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, mjpb)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) JobPoolRemoveHandler(request *restful.Request, response *restful.Response) {

	jpbean := new(module.MetaJobPoolBean)
	err := request.ReadEntity(&jpbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jpbean.FlowId) == 0 || len(jpbean.Sys) == 0 || len(jpbean.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("flowid sys or job missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flowid sys job missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_JOBSPOOL)
	defer bt.Close()

	err = bt.Remove(jpbean.FlowId + "." + jpbean.Sys + "." + jpbean.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemRingGoListHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	for k := range ringGoSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingGoOffsetBean)
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemRingGoRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRingGoRemoveBean)
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
	ringGoSpool.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemRingPendingListHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	for k := range ringPendingSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceSystem) SystemRingPendingRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRingPendingRemoveBean)
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
	ringPendingSpool.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}
