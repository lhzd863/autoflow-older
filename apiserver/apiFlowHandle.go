package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/db"
	"github.com/lhzd863/autoflow/glog"
	"github.com/lhzd863/autoflow/gproto"
	"github.com/lhzd863/autoflow/module"
	"github.com/lhzd863/autoflow/util"
)

type ResponseResourceFlow struct {
	sync.Mutex
}

func NewResponseResourceFlow() *ResponseResourceFlow {
	return &ResponseResourceFlow{}
}

func (rrs *ResponseResourceFlow) InstanceCreateHandler(request *restful.Request, response *restful.Response) {
	reqParams, err := url.ParseQuery(request.Request.URL.RawQuery)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse para error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse para error.%v", err), nil)
		return
	}

	username, err := util.JwtAccessTokenUserName(fmt.Sprint(reqParams["accesstoken"][0]), conf.JwtKey)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get username error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get username error.%v", err), nil)
		return
	}
	mjf := new(module.MetaParaInstanceCreateBean)
	err = request.ReadEntity(&mjf)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	pn, err := strconv.Atoi(mjf.ProcessNum)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("para conv processnum int error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("para conv processnum int error.%v", err), nil)
		return
	}
	mparr := rrs.getMstPort(pn)
	if len(mparr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no mst start flow."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no mst start flow."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()

	ib := bt.Get(mjf.ImageId)
	bt.Close()
	if ib == nil {
		glog.Glog(LogF, fmt.Sprintf("%v image no data result.", mjf.ImageId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v image no data result.", mjf.ImageId), nil)
		return
	}
	masb := new(module.MetaJobImageBean)
	err = json.Unmarshal([]byte(ib.(string)), &masb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	flowinstbean := new(module.MetaJobFlowBean)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	flowinstbean.CreateTime = timeStr
	flowinstbean.RunContext = timeStr
	flowinstbean.StartTime = timeStr
	flowinstbean.ImageId = masb.ImageId
	flowinstbean.Enable = masb.Enable
	flowinstbean.User = username
	flowinstbean.Status = util.STATUS_AUTO_RUNNING
	flowinstbean.HomeDir = conf.HomeDir

	frid := util.RandStringBytes(6)
	timeStrId := time.Now().Format("20060102150405")
	fid := fmt.Sprintf("%v%v", timeStrId, frid)
	flowinstbean.FlowId = fid
	flowinstbean.DbStore = fmt.Sprintf("%v/%v/%v.db", conf.HomeDir, flowinstbean.FlowId, flowinstbean.FlowId)

	if ok, err := util.PathExists(conf.HomeDir + "/" + flowinstbean.FlowId); !ok {
		os.Mkdir(conf.HomeDir+"/"+flowinstbean.FlowId, os.ModePerm)
	} else {
		glog.Glog(LogF, fmt.Sprint(err))
	}

	if !util.FileExist(conf.HomeDir + "/" + flowinstbean.FlowId + "/" + flowinstbean.FlowId + ".db") {
		util.Copy(masb.DbStore, conf.HomeDir+"/"+flowinstbean.FlowId+"/"+flowinstbean.FlowId+".db")
	}
	bt1 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt1.Close()

	jsonstr, err := json.Marshal(flowinstbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal error.%v", err), nil)
		return
	}
	err = bt1.Set(flowinstbean.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt1.Close()
	for _, mh := range mparr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(mh.Ip+":"+mh.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		mifb.ProcessNum = mh.ProcessNum
		mifb.FlowId = fid

		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowStart(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}

		mmfb := new(module.MetaMstFlowBean)
		mmfb.FlowId = fid
		mmfb.MstId = mh.MstId
		mmfb.Ip = mh.Ip
		mmfb.Port = mh.Port
		mmfb.UpdateTime = timeStr
		mmfb.ProcessNum = flowinstbean.ProcessNum
		err = rrs.saveMstFlow(mmfb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not save: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not save: %v", mh.Ip, mh.Port, err), nil)
			return
		}
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) InstanceStartHandler(request *restful.Request, response *restful.Response) {
	mjfbean := new(module.MetaInstanceFlowBean)
	err := request.ReadEntity(&mjfbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	pn, err := strconv.Atoi(mjfbean.ProcessNum)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("para conv processnum int error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("para conv processnum int error.%v", err), nil)
		return
	}

	mparr := rrs.getMstPort(pn)
	if len(mparr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no mst start flow."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no mst start flow."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(mjfbean.FlowId)
	if fb0 == nil {
		glog.Glog(LogF, fmt.Sprintf("flow %v not exists,start flow fail.", mjfbean.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flow %v not exists,start flow fail.", mjfbean.FlowId), nil)
		return
	}
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if fb == nil {
		glog.Glog(LogF, fmt.Sprintf("%v no data result", fb.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v no data result", fb.FlowId), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.StartTime = timeStr
	fb.EndTime = ""
	fb.Status = util.STATUS_AUTO_RUNNING

	if !util.FileExist(fb.HomeDir + "/" + fb.FlowId + "/" + fb.FlowId + ".db") {
		glog.Glog(LogF, fb.HomeDir+"/"+fb.FlowId+"/"+fb.FlowId+".db not exists.")
		util.ApiResponse(response.ResponseWriter, 700, fb.HomeDir+"/"+fb.FlowId+"/"+fb.FlowId+".db not exists.", nil)
		return
	}

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt.Close()
	for _, mh := range mparr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(mh.Ip+":"+mh.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		mifb.ProcessNum = mjfbean.ProcessNum
		mifb.FlowId = mjfbean.FlowId

		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowStart(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}

		mmfb := new(module.MetaMstFlowBean)
		mmfb.FlowId = fb.FlowId
		mmfb.MstId = mh.MstId
		mmfb.Ip = mh.Ip
		mmfb.Port = mh.Port
		mmfb.UpdateTime = timeStr
		mmfb.ProcessNum = mjfbean.ProcessNum
		err = rrs.saveMstFlow(mmfb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not save: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not save: %v", mh.Ip, mh.Port, err), nil)
			return
		}
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowListHandler(request *restful.Request, response *restful.Response) {
	mharr := rrs.getMstHeart()
	if len(mharr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no flow port error."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no flow port error."), nil)
		return
	}
	tlst := make([]interface{}, 0)
	for _, mh := range mharr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(mh.Ip+":"+mh.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowStatus(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}
		if len(tr.Data) == 0 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			break
		}
		mf := new([]module.MetaFlowStatusBean)
		err = json.Unmarshal([]byte(tr.Data), &mf)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json err.%v %v", err, tr.Data), nil)
			return
		}
		for _, v := range *mf {
			tlst = append(tlst, v)
		}
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mjfb := new(module.MetaJobFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjfb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			mjfb.Status = util.STATUS_AUTO_STOP
			for _, tv := range tlst {
				tv1 := tv.(module.MetaFlowStatusBean)
				if tv1.FlowId == mjfb.FlowId {
					mjfb.Status = util.STATUS_AUTO_RUNNING
				}
			}
			retlst = append(retlst, mjfb)
		}
	}

	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_JOB)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.FlowId)
	if ib != nil {
		m := new(module.MetaJobFlowBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(p.FlowId)
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.DbStore = p.DbStore
	fb.HomeDir = p.HomeDir
	fb.RunContext = p.RunContext
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) MstInstanceStartHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowStartBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	pn, err := strconv.Atoi(p.ProcessNum)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("para conv processnum int error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("para conv processnum int error.%v", err), nil)
		return
	}

	m, err := rrs.getMstInfo(p.MstId, pn)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get mst info err.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get mst info err.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(p.FlowId)
	if fb0 == nil {
		glog.Glog(LogF, fmt.Sprintf("flow %v not exists,start flow fail.", p.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flow %v not exists,start flow fail.", p.FlowId), nil)
		return
	}
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if fb == nil {
		glog.Glog(LogF, fmt.Sprintf("%v no data result", fb.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v no data result", fb.FlowId), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.StartTime = timeStr
	fb.EndTime = ""
	fb.Status = util.STATUS_AUTO_RUNNING

	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}

	if !util.FileExist(dbf) {
		glog.Glog(LogF, dbf+" not exists.")
		util.ApiResponse(response.ResponseWriter, 700, dbf+" not exists.", nil)
		return
	}

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt.Close()

	// 建立连接到gRPC服务
	conn, err := grpc.Dial(m.Ip+":"+m.Port, grpc.WithInsecure())
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", m.Ip, m.Port, err), nil)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := gproto.NewFlowMasterClient(conn)

	mifb := new(module.MetaInstanceFlowBean)
	mifb.ProcessNum = m.ProcessNum
	mifb.FlowId = p.FlowId

	jsonstr, err = json.Marshal(mifb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
		return
	}
	// 调用gRPC接口
	tr, err := t.FlowStart(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", m.Ip, m.Port, err), nil)
		return
	}
	//change job status
	if tr.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
		return
	}

	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) InstanceStopHandler(request *restful.Request, response *restful.Response) {
	flowbean := new(module.MetaJobFlowBean)
	err := request.ReadEntity(&flowbean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	fparr, err := rrs.getFlowPort(flowbean.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow port error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get mst port error.%v", err), nil)
		return
	}

	for _, fp := range fparr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(fp.Ip+":"+fp.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		mifb.FlowId = fp.FlowId

		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowStop(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(flowbean.FlowId)
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Status = util.STATUS_AUTO_STOP
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.EndTime = timeStr
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt.Close()
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) MstInstanceStopHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowStopBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	m, err := rrs.getMstInfo(p.MstId, 1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get mst info error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get mst info error.%v", err), nil)
		return
	}

	// 建立连接到gRPC服务
	conn, err := grpc.Dial(m.Ip+":"+m.Port, grpc.WithInsecure())
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", m.Ip, m.Port, err), nil)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := gproto.NewFlowMasterClient(conn)

	mifb := new(module.MetaInstanceFlowBean)
	mifb.FlowId = p.FlowId

	jsonstr, err := json.Marshal(mifb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
		return
	}
	// 调用gRPC接口
	tr, err := t.FlowStop(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", m.Ip, m.Port, err), nil)
		return
	}
	//change job status
	if tr.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(p.FlowId)
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Status = util.STATUS_AUTO_STOP
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.EndTime = timeStr
	jsonstr, _ = json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt.Close()
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) MstFlowRoutineStartHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowStartBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	m, err := rrs.getMstInfo(p.MstId, 1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get mst info err.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get mst info err.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()
	fb0 := bt.Get(p.FlowId)
	if fb0 == nil {
		glog.Glog(LogF, fmt.Sprintf("flow %v not exists,start flow fail.", p.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flow %v not exists,start flow fail.", p.FlowId), nil)
		return
	}
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if fb == nil {
		glog.Glog(LogF, fmt.Sprintf("%v no data result", fb.FlowId))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v no data result", fb.FlowId), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.StartTime = timeStr
	fb.EndTime = ""
	fb.Status = util.STATUS_AUTO_RUNNING

	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}

	if !util.FileExist(dbf) {
		glog.Glog(LogF, dbf+" not exists.")
		util.ApiResponse(response.ResponseWriter, 700, dbf+" not exists.", nil)
		return
	}

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	bt.Close()

	// 建立连接到gRPC服务
	conn, err := grpc.Dial(m.Ip+":"+m.Port, grpc.WithInsecure())
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", m.Ip, m.Port, err), nil)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := gproto.NewFlowMasterClient(conn)

	mifb := new(module.MetaInstanceFlowBean)
	mifb.ProcessNum = "1"
	mifb.FlowId = p.FlowId
	mifb.RoutineId = p.RoutineId

	jsonstr, err = json.Marshal(mifb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
		return
	}
	// 调用gRPC接口
	tr, err := t.MstFlowRoutineStart(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", m.Ip, m.Port, err), nil)
		return
	}
	//change job status
	if tr.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
		return
	}

	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) MstFlowRoutineStopHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowStartBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	m, err := rrs.getMstInfo(p.MstId, 1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get mst info err.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get mst info err.%v", err), nil)
		return
	}

	// 建立连接到gRPC服务
	conn, err := grpc.Dial(m.Ip+":"+m.Port, grpc.WithInsecure())
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", m.Ip, m.Port, err), nil)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := gproto.NewFlowMasterClient(conn)

	mifb := new(module.MetaInstanceFlowBean)
	mifb.ProcessNum = "1"
	mifb.FlowId = p.FlowId
	mifb.RoutineId = p.RoutineId

	jsonstr, err := json.Marshal(mifb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
		return
	}
	// 调用gRPC接口
	tr, err := t.MstFlowRoutineStop(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", m.Ip, m.Port, err), nil)
		return
	}
	//change job status
	if tr.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
		return
	}

	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) InstanceListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mjfb := new(module.MetaJobFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjfb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, mjfb)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) InstanceListStatusHandler(request *restful.Request, response *restful.Response) {

	bt0 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
	defer bt0.Close()

	strlist0 := bt0.Scan()
	retlst0 := make([]*module.MetaMstFlowRoutineHeartBean, 0)
	for _, v := range strlist0 {
		for k1, v1 := range v.(map[string]interface{}) {
			mmhb := new(module.MetaMstFlowRoutineHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmhb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(mmhb.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", mmhb.MstId, mmhb.Ip, mmhb.Port))
				bt0.Remove(k1)
				continue
			}
			retlst0 = append(retlst0, mmhb)
		}
	}
	bt0.Close()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mjfb := new(module.MetaJobFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjfb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			mjfb.Status = util.STATUS_AUTO_STOP
			for _, tv := range retlst0 {
				tv1 := tv
				if tv1.FlowId == mjfb.FlowId {
					mjfb.Status = util.STATUS_AUTO_RUNNING
				}
			}
			retlst = append(retlst, mjfb)
		}
	}

	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) flowInfo(flowid string) (*module.MetaJobFlowBean, error) {
	rrs.Lock()
	defer rrs.Unlock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	var f *module.MetaJobFlowBean
	flag := 0
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if m.FlowId != flowid {
				continue
			}
			f = m
			flag = 1
		}
	}
	if flag == 0 {
		return f, errors.New(fmt.Sprintf("%v no flow information.", flowid))
	}
	return f, nil
}

func (rrs *ResponseResourceFlow) InstanceRemoveHandler(request *restful.Request, response *restful.Response) {
	m := new(module.MetaJobFlowBean)
	err := request.ReadEntity(&m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW)
	defer bt.Close()

	err = bt.Remove(m.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) MstFlowRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaMstFlowRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
	defer bt.Close()

	err = bt.Remove(p.MstId + "," + p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowRoutineStatusHandler(request *restful.Request, response *restful.Response) {
	mharr := rrs.getMstHeart()
	if len(mharr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no flow port error."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no flow port error."), nil)
		return
	}
	retlst := make([]interface{}, 0)
	for _, mh := range mharr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(mh.Ip+":"+mh.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowStatus(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", mh.Ip, mh.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}
		if len(tr.Data) == 0 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 200, "", retlst)
			return
		}
		mf := new([]module.MetaFlowStatusBean)
		err = json.Unmarshal([]byte(tr.Data), &mf)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json err.%v %v", err, tr.Data), nil)
			return
		}
		for _, v := range *mf {
			retlst = append(retlst, v)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowStatusUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowStatusUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	bt2 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_JOB)
	defer bt2.Close()
	fb0 := bt2.Get(p.FlowId)
	fb := new(module.MetaJobFlowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Status = p.Status
	jsonstr, _ := json.Marshal(fb)
	err = bt2.Set(fb.FlowId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowRoutineAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowRoutineAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	glog.Glog(LogF, fmt.Sprint(p))
	fparr := rrs.getMstFlowPort(p.FlowId, p.MstId)
	if len(fparr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no mst flow routine error."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no mst flow routine error."), nil)
		return
	}

	for _, fp := range fparr {
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(fp.Ip+":"+fp.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		mifb.ProcessNum = p.ProcessNum
		mifb.FlowId = fp.FlowId

		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowRoutineAdd(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowRoutineSubHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowRoutineSubBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	fparr := rrs.getMstFlowPort(p.FlowId, p.MstId)
	if len(fparr) == 0 {
		glog.Glog(LogF, fmt.Sprintf("no mst flow routine error."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("no mst flow routine error."), nil)
		return
	}

	for _, fp := range fparr {
		glog.Glog(LogF, fmt.Sprintf("flow %v on %v:%v stop %v routine.", p.FlowId, fp.Ip, fp.Port, p.ProcessNum))
		// 建立连接到gRPC服务
		conn, err := grpc.Dial(fp.Ip+":"+fp.Port, grpc.WithInsecure())
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v did not connect: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 创建Waiter服务的客户端
		t := gproto.NewFlowMasterClient(conn)

		mifb := new(module.MetaInstanceFlowBean)
		mifb.ProcessNum = p.ProcessNum
		mifb.FlowId = fp.FlowId

		jsonstr, err := json.Marshal(mifb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("json marshal %v", err), nil)
			return
		}
		// 调用gRPC接口
		tr, err := t.FlowRoutineSub(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v:%v could not greet: %v", fp.Ip, fp.Port, err), nil)
			return
		}
		//change job status
		if tr.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprint(tr.Status_Txt), nil)
			return
		}
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowParameterListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
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

func (rrs *ResponseResourceFlow) FlowParameterGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowParameterGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.FlowId + "." + p.Key)
	if ib != nil {
		m := new(module.MetaJobParameterBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowParameterUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowParameterUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt.Close()
	fb0 := bt.Get(p.FlowId + "." + p.Key)
	fb := new(module.MetaJobParameterBean)
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
	err = bt.Set(fb.FlowId+"."+fb.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowParameterRemoveHandler(request *restful.Request, response *restful.Response) {
	m := new(module.MetaParaFlowParameterRemoveBean)
	err := request.ReadEntity(&m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt.Close()

	err = bt.Remove(m.FlowId + "." + m.Key)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) FlowParameterAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowParameterAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}

	if len(p.FlowId) == 0 || len(p.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("flowid or key missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flowid or key missed."), nil)
		return
	}
	p.Type = "F"
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt.Close()

	jsonstr, _ := json.Marshal(p)
	err = bt.Set(p.FlowId+"."+p.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceFlow) getMstPort(pnum int) []*module.MetaMstHeartBean {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]*module.MetaMstHeartBean, 0)
	cnt := 1
	tmap := make(map[string]*module.MetaMstHeartBean)
	for j := 0; j < pnum && cnt <= pnum; j++ {
		for _, v := range strlist {
			if cnt > pnum {
				break
			}
			for _, v1 := range v.(map[string]interface{}) {
				if cnt > pnum {
					break
				}
				mmhb := new(module.MetaMstHeartBean)
				err := json.Unmarshal([]byte(v1.(string)), &mmhb)
				if err != nil {
					glog.Glog(LogF, fmt.Sprint(err))
					continue
				}
				timeStr := time.Now().Format("2006-01-02 15:04:05")
				ise, _ := util.IsExpired(mmhb.UpdateTime, timeStr, 300)
				if ise {
					glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", mmhb.MstId, mmhb.Ip, mmhb.Port))
					continue
				}
				_, ok := tmap[mmhb.MstId]
				if ok {
					v2 := tmap[mmhb.MstId]
					pn, _ := strconv.Atoi(v2.ProcessNum)
					v2.ProcessNum = fmt.Sprint(pn + 1)
					tmap[v2.MstId] = v2
				} else {
					mmhb.ProcessNum = fmt.Sprint(1)
					tmap[mmhb.MstId] = mmhb
				}
				cnt += 1
			}
		}
	}
	for k0 := range tmap {
		mmh := tmap[k0]
		retlst = append(retlst, mmh)
	}
	return retlst
}

func (rrs *ResponseResourceFlow) saveMstFlow(mmfb *module.MetaMstFlowBean) error {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
	defer bt.Close()

	jsonstr, _ := json.Marshal(mmfb)
	err := bt.Set(mmfb.MstId+"."+mmfb.FlowId, string(jsonstr))
	if err != nil {
		return errors.New(fmt.Sprintf("data in db update error.%v", err))
	}
	return nil
}

func (rrs *ResponseResourceFlow) getMstHeart() []*module.MetaMstHeartBean {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]*module.MetaMstHeartBean, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaMstHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", m.MstId, m.Ip, m.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, m)
		}
	}
	return retlst
}

func (rrs *ResponseResourceFlow) getMstInfo(mstid string, pnum int) (*module.MetaMstHeartBean, error) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	flag := 0
	var ret *module.MetaMstHeartBean
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaMstHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(m.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v:%v.", m.MstId, m.Ip, m.Port))
				continue
			}
			if m.MstId != mstid {
				continue
			}
			ret = m
			ret.ProcessNum = fmt.Sprint(pnum)
			flag = 1
		}
	}
	for flag == 0 {
		return nil, errors.New("no mst obtain.")
	}
	return ret, nil
}

func (rrs *ResponseResourceFlow) getFlowPort(flowid string) ([]module.MetaMstFlowBean, error) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	flag := 0
	retlst := make([]module.MetaMstFlowBean, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			mmf := new(module.MetaMstFlowBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmf)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if mmf.FlowId != flowid {
				glog.Glog(LogF, fmt.Sprintf("%v != %v", mmf.FlowId, flowid))
				continue
			}
			ism, _ := rrs.IsExpiredMst(mmf.MstId)
			if ism {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", mmf.MstId, mmf.Ip, mmf.Port))
				bt0 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
				defer bt0.Close()
				bt0.Remove(k1)
				bt0.Close()
				continue
			}
			retlst = append(retlst, *mmf)
			flag = 1
		}
	}

	if flag == 0 {
		return nil, errors.New(fmt.Sprintf("flow %v no ip and port,wait for next time.", flowid))
	}
	return retlst, nil
}

func (rrs *ResponseResourceFlow) getMstFlowPort(flowid string, mstid string) []*module.MetaMstFlowRoutineHeartBean {
	rrs.Lock()
	rrs.Unlock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_ROUTINE_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]*module.MetaMstFlowRoutineHeartBean, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			mmf := new(module.MetaMstFlowRoutineHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmf)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if mmf.FlowId != flowid || mmf.MstId != mstid {
				glog.Glog(LogF, fmt.Sprintf("%v!=%v or %v!=%v.", mmf.FlowId, flowid, mmf.MstId, mstid))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(mmf.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v:%v.", mmf.MstId, mmf.Ip, mmf.Port))
				bt.Remove(k1)
				continue
			}
			retlst = append(retlst, mmf)
			return retlst
		}
	}
	return retlst
}

func (rrs *ResponseResourceFlow) IsExpiredMst(mstid string) (bool, error) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	flag := 1
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mmh := new(module.MetaMstHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmh)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if mmh.MstId != mstid {
				glog.Glog(LogF, fmt.Sprintf("%v!=%v.", mmh.MstId, mstid))
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(mmh.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v:%v.", mmh.MstId, mmh.Ip, mmh.Port))
				flag = 1
				continue
			}
			flag = 0
		}
	}
	if flag == 1 {
		return true, nil
	}
	return false, nil
}

func (rrs *ResponseResourceFlow) instanceHomeDir(flowid string) (string, error) {

	return conf.HomeDir, nil
}

func (rrs *ResponseResourceFlow) flowDbFile(flowid string) (string, error) {
	f := conf.HomeDir + "/" + flowid + "/" + flowid + ".db"
	return f, nil
}
