package mst

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/satori/go.uuid"

	"github.com/lhzd863/autoflow/db"
	"github.com/lhzd863/autoflow/glog"
	"github.com/lhzd863/autoflow/gproto"
	"github.com/lhzd863/autoflow/module"
	"github.com/lhzd863/autoflow/util"
	"github.com/lhzd863/autoflow/workpool"
)

var (
	LogF                string
	HomeDir             string
	AccessToken         string
	ApiServerIp         string
	ApiServerPort       string
	flowspool           = db.NewMemDB()
	jobpool             = db.NewMemDB()
	workerpool          = db.NewMemDB()
	ProcessNum          int
	Ip                  string
	Port                string
	MstId               string
	shutdownFlowChannel = make(chan interface{})
	flowChannel         = make(chan interface{})
)

// 业务实现方法的容器
type MServer struct {
	sync.RWMutex
	waitGroup util.WaitGroupWrapper
}

func NewMServer(paraMap map[string]interface{}) *MServer {
	MstId = paraMap["mstid"].(string)
	Ip = paraMap["ip"].(string)
	Port = paraMap["port"].(string)
	HomeDir = paraMap["homedir"].(string)
	AccessToken = paraMap["accesstoken"].(string)
	ApiServerIp = paraMap["apiserverip"].(string)
	ApiServerPort = paraMap["apiserverport"].(string)
	NumStr := paraMap["processnum"].(string)
	ProcessNum, _ = strconv.Atoi(NumStr)
	LogF = HomeDir + "/mst_${" + util.ENV_VAR_DATE + "}.log"
	m := &MServer{}
	return m

}

func (ms *MServer) DoCmd(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) Ping(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) JobStart(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) JobExecuter(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) FlowStart(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mifb := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mifb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	if flowspool.Get(mifb.FlowId) != nil && len(flowspool.Get(mifb.FlowId).(map[string]interface{})) != 0 {
		glog.Glog(LogF, fmt.Sprintf("%v has exists.", mifb.FlowId))
		return &gproto.Res{Status_Txt: fmt.Sprintf("%v has exists.", mifb.FlowId), Status_Code: 700, Data: "{}"}, nil
	}

	pnum, err := strconv.Atoi(mifb.ProcessNum)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	workPool := workpool.New(pnum, 1000)
	mp := make(map[string]interface{})
	wp := make([]interface{}, 0)
	mp["workpool"] = workPool
	cnt := 0
	glog.Glog(LogF, fmt.Sprintf("flowid %v run %v %v", mifb.FlowId))
	for j := 0; j < pnum; j++ {
		mw := new(MetaMyWork)
		mw.FlowId = mifb.FlowId
		mw.MstIp = Ip
		mw.MstPort = Port
		mw.ApiServerIp = ApiServerIp
		mw.ApiServerPort = ApiServerPort
		mw.HomeDir = HomeDir
		mw.AccessToken = AccessToken
		work := &MyWork{
			Id:     fmt.Sprint(j),
			Name:   "A-" + fmt.Sprint(j),
			MstId:  MstId,
			FlowId: mifb.FlowId,
			WP:     workPool,
			Mmw:    mw,
		}
		err := workPool.PostWork(fmt.Sprintf("name_routine_%v", j), work)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("ERROR: %s", err))
			continue
		}
		wp = append(wp, work)
		cnt++
		glog.Glog(LogF, fmt.Sprintf("new add routine A-%v", j))
	}
	glog.Glog(LogF, fmt.Sprintf("total new add %v routine.", cnt))
	mp["mywork"] = wp
	flowspool.Add(mifb.FlowId, mp)
	if cnt == pnum {
		return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
	}
	msg := fmt.Sprintf("total %v proccess need start, current %v start ,not all process start.", pnum, cnt)
	return &gproto.Res{Status_Txt: msg, Status_Code: 700, Data: "{}"}, nil
}

func (ms *MServer) FlowCreate(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 700, Data: "{}"}, nil
}

func (ms *MServer) FlowStop(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mifb := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mifb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
	}
	if flowspool.Get(mifb.FlowId) != nil && len(flowspool.Get(mifb.FlowId).(map[string]interface{})) != 0 {
		glog.Glog(LogF, fmt.Sprintf("flow %v will stop.", mifb.FlowId))
		mp := (flowspool.Get(mifb.FlowId)).(map[string]interface{})
		wp := mp["workpool"].(*workpool.WorkPool)
		mw := mp["mywork"].([]interface{})
		cnt := 0
		for k1 := range mw {
			wk1 := (mw[k1]).(*MyWork)
			wk1.StopFlag = "1"
			glog.Glog(LogF, fmt.Sprintf("flow %v stop routine %v.", mifb.FlowId, wk1.Name))
			cnt++
		}
		wp.Shutdown("A-" + fmt.Sprint(0))
		glog.Glog(LogF, fmt.Sprintf("total remove %v routine.", cnt))
		flowspool.Remove(mifb.FlowId)
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) FlowStatus(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	retlst := make([]interface{}, 0)
	for k := range flowspool.MemMap {
		mfsb := new(module.MetaFlowStatusBean)
		mfsb.WorkPoolStatus = ""
		mfsb.MyWorkCnt = "0"
		mp := (flowspool.MemMap[k]).(map[string]interface{})
		if mp == nil {
			continue
		}
		wp := mp["workpool"].(*workpool.WorkPool)
		mw := mp["mywork"].([]interface{})
		if wp != nil {
			mfsb.WorkPoolStatus = "Running"
		}
		if mw != nil {
			mfsb.MyWorkCnt = fmt.Sprint(len(mw))
		}
		mfsb.Ip = Ip
		mfsb.Port = Port
		mfsb.MstId = MstId
		mfsb.FlowId = k
		retlst = append(retlst, mfsb)
	}
	jsonstr, err := json.Marshal(retlst)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		return &gproto.Res{Status_Txt: "stataus conv string error", Status_Code: 700, Data: "{}"}, nil
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: string(jsonstr)}, nil
}

func (ms *MServer) FlowRoutineAdd(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mifb := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mifb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	if flowspool.Get(mifb.FlowId) == nil || len(flowspool.Get(mifb.FlowId).(map[string]interface{})) == 0 {
		glog.Glog(LogF, fmt.Sprintf("%v no exists.", mifb.FlowId))
		return &gproto.Res{Status_Txt: fmt.Sprintf("%v no exists.", mifb.FlowId), Status_Code: 700, Data: "{}"}, nil
	}
	mp := (flowspool.Get(mifb.FlowId)).(map[string]interface{})
	wp := mp["workpool"].(*workpool.WorkPool)
	mw := mp["mywork"].([]interface{})
	pnum, err := strconv.Atoi(mifb.ProcessNum)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	totalnum := int(wp.QueuedWork()+wp.ActiveRoutines()) + pnum
	arr := make([]string, 0)
	for i := 0; i < totalnum; i++ {
		flg := 0
		for k := range mw {
			wk0 := mw[k].(*MyWork)
			if fmt.Sprint(i) == wk0.Id {
				flg = 1
				break
			}
		}
		if flg == 0 {
			arr = append(arr, fmt.Sprint(i))
		}
		if len(arr) >= pnum {
			break
		}
	}
	cnt := 0
	glog.Glog(LogF, fmt.Sprintf("%v need run runtine %v.", mifb.FlowId, len(arr)))
	for a1 := range arr {
		ii := arr[a1]
		ri, _ := strconv.Atoi(ii)
		wp.AddRoutine(ri)
		mw1 := new(MetaMyWork)
		mw1.FlowId = mifb.FlowId
		mw1.MstIp = Ip
		mw1.MstPort = Port
		mw1.ApiServerIp = ApiServerIp
		mw1.ApiServerPort = ApiServerPort
		mw1.HomeDir = HomeDir
		mw1.AccessToken = AccessToken
		work := &MyWork{
			Id:     ii,
			Name:   "A-" + ii,
			MstId:  MstId,
			FlowId: mifb.FlowId,
			WP:     wp,
			Mmw:    mw1,
		}
		err = wp.PostWork(fmt.Sprintf("name_routine_%v", ii), work)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("ERROR: %s\n", err))
			continue
		}
		mw = append(mw, work)
		cnt++
		glog.Glog(LogF, fmt.Sprintf("new add routine A-%v", ii))
	}
	mp["mywork"] = mw
	mp["workpool"] = wp
	flowspool.Add(mifb.FlowId, mp)
	glog.Glog(LogF, fmt.Sprintf("total add %v routine.", cnt))
	if cnt == pnum {
		return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
	}
	msg := fmt.Sprintf("total %v proccess need start, current %v start ,not all process start.", pnum, cnt)
	return &gproto.Res{Status_Txt: msg, Status_Code: 700, Data: "{}"}, nil
}

func (ms *MServer) FlowRoutineSub(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mifb := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mifb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	pnum, err := strconv.Atoi(mifb.ProcessNum)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	if flowspool.Get(mifb.FlowId) != nil && len(flowspool.Get(mifb.FlowId).(map[string]interface{})) != 0 {
		glog.Glog(LogF, fmt.Sprintf("flow %v will stop.", mifb.FlowId))
		mp := (flowspool.Get(mifb.FlowId)).(map[string]interface{})
		mw := mp["mywork"].([]interface{})
		cnt := 0
		nmw := make([]interface{}, 0)
		for k1 := range mw {
			if cnt >= pnum {
				nmw = append(nmw, mw[k1])
				continue
			}
			wk1 := (mw[k1]).(*MyWork)
			wk1.StopFlag = "1"
			glog.Glog(LogF, fmt.Sprintf("flow %v stop routine %v.", mifb.FlowId, wk1.Name))
			cnt++
		}
		mp["mywork"] = nmw
		flowspool.Add(mifb.FlowId, mp)
		glog.Glog(LogF, fmt.Sprintf("total remove %v routine.", cnt))
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) FlowRoutineList(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mifb := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mifb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	mfsb := new(module.MetaFlowStatusBean)
	mfsb.WorkPoolStatus = ""
	mfsb.MyWorkCnt = "0"
	for k := range flowspool.MemMap {
		if k != mifb.FlowId {
			continue
		}
		mp := (flowspool.MemMap[k]).(map[string]interface{})
		if mp == nil {
			continue
		}
		wp := mp["workpool"].(*workpool.WorkPool)
		mw := mp["mywork"].([]interface{})
		mfsb.Ip = Ip
		mfsb.Port = Port
		mfsb.FlowId = k
		mfsb.MstId = MstId
		if wp != nil {
			mfsb.WorkPoolStatus = "Running"
		}
		if mw != nil {
			mfsb.MyWorkCnt = fmt.Sprint(len(mw))
		}
	}
	jsonstr, err := json.Marshal(mfsb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		return &gproto.Res{Status_Txt: "stataus conv string error", Status_Code: 700, Data: "{}"}, nil
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: string(jsonstr)}, nil
}

func (ms *MServer) JobStop(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) JobStatus(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) MstFlowRoutineStop(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	p := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &p)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	if flowspool.Get(p.FlowId) != nil && len(flowspool.Get(p.FlowId).(map[string]interface{})) != 0 {
		glog.Glog(LogF, fmt.Sprintf("flow %v will stop.", p.FlowId))
		mp := (flowspool.Get(p.FlowId)).(map[string]interface{})
		mw := mp["mywork"].([]interface{})
		nmw := make([]interface{}, 0)
		for k1 := range mw {
			wk1 := (mw[k1]).(*MyWork)
			if wk1.Id == p.RoutineId {
				wk1.StopFlag = "1"
				glog.Glog(LogF, fmt.Sprintf("flow %v stop routine %v.", p.MstId, p.FlowId, wk1.Id))
				continue
			}
			nmw = append(nmw, mw[k1])
		}
		mp["mywork"] = nmw
		flowspool.Add(p.FlowId, mp)
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) MstFlowRoutineStart(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	p := new(module.MetaInstanceFlowBean)
	err := json.Unmarshal([]byte(in.JsonStr), &p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("%v", err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	if flowspool.Get(p.FlowId) == nil || len(flowspool.Get(p.FlowId).(map[string]interface{})) == 0 {
		glog.Glog(LogF, fmt.Sprintf("%v no exists.", p.FlowId))
		return &gproto.Res{Status_Txt: fmt.Sprintf("%v no exists.", p.FlowId), Status_Code: 700, Data: "{}"}, nil
	}
	mp := (flowspool.Get(p.FlowId)).(map[string]interface{})
	wp := mp["workpool"].(*workpool.WorkPool)
	mw := mp["mywork"].([]interface{})

	for k := range mw {
		wk0 := mw[k].(*MyWork)
		if p.RoutineId == wk0.Id {
			glog.Glog(LogF, fmt.Sprintf("%v mst %v flow %v routine has exists.", p.MstId, p.FlowId, p.RoutineId))
			return &gproto.Res{Status_Txt: fmt.Sprintf("%v mst %v flow %v routine has exists.", p.MstId, p.FlowId, p.RoutineId), Status_Code: 700, Data: "{}"}, nil
		}
	}

	glog.Glog(LogF, fmt.Sprintf("%v %v need run runtine %v.", p.MstId, p.FlowId, p.RoutineId))

	ri, _ := strconv.Atoi(p.RoutineId)
	wp.AddRoutine(ri)
	mw1 := new(MetaMyWork)
	mw1.FlowId = p.FlowId
	mw1.MstIp = Ip
	mw1.MstPort = Port
	mw1.ApiServerIp = ApiServerIp
	mw1.ApiServerPort = ApiServerPort
	mw1.HomeDir = HomeDir
	mw1.AccessToken = AccessToken
	work := &MyWork{
		Id:     p.RoutineId,
		Name:   "A-" + p.RoutineId,
		MstId:  MstId,
		FlowId: p.FlowId,
		WP:     wp,
		Mmw:    mw1,
	}
	err = wp.PostWork(fmt.Sprintf("name_routine_%v", p.RoutineId), work)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("ERROR: %s\n", err))
		return &gproto.Res{Status_Txt: fmt.Sprint(err), Status_Code: 700, Data: "{}"}, nil
	}
	mw = append(mw, work)

	glog.Glog(LogF, fmt.Sprintf("new add routine A-%v", p.RoutineId))

	mp["mywork"] = mw
	mp["workpool"] = wp
	flowspool.Add(p.FlowId, mp)

	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (ms *MServer) Main() bool {
	var wg util.WaitGroupWrapper
	StopFlag := 0

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	m := new(module.MetaParaMstHeartAddBean)
	u1 := uuid.Must(uuid.NewV4(),nil)
	m.Id = fmt.Sprint(u1)
	m.MstId = MstId
	m.Ip = Ip
	m.Port = Port
	timeStr0 := time.Now().Format("2006-01-02 15:04:05")
	m.StartTime = timeStr0

	go func() {
		<-signalChan
		StopFlag = 1
		exitChan <- 1
		ms.RegisterRemove(m)
	}()
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	wg.Wrap(func() {
		st := time.Now().Unix()
		et := time.Now().Unix()
		for {
			if StopFlag == 1 {
				break
			}
			et = time.Now().Unix()
			if et-st > 30 {
				ret := ms.Register(m)
				if !ret {
					glog.Glog(LogF, "register mst fail.")
				}
				st = time.Now().Unix()
			}
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(2)
			time.Sleep(time.Duration(ri) * time.Second)
		}
	})

	lis, err := net.Listen("tcp", ":"+Port) //监听所有网卡8028端口的TCP连接
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("监听失败: %v", err))
		return false
	}
	ss := grpc.NewServer() //创建gRPC服务

	/**注册接口服务
	 * 以定义proto时的service为单位注册，服务中可以有多个方法
	 * (proto编译时会为每个service生成Register***Server方法)
	 * 包.注册服务方法(gRpc服务实例，包含接口方法的结构体[指针])
	 */
	gproto.RegisterFlowMasterServer(ss, &MServer{})
	/**如果有可以注册多个接口服务,结构体要实现对应的接口方法
	 * user.RegisterLoginServer(s, &server{})
	 * minMovie.RegisterFbiServer(s, &server{})
	 */
	// 在gRPC服务器上注册反射服务
	reflection.Register(ss)
	// 将监听交给gRPC服务处理
	wg.Wrap(func() {
		err = ss.Serve(lis)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("failed to serve: %v", err))
			return
		}
	})
	<-exitChan
	StopFlag = 1
	ss.Stop()
	wg.Wait()
	glog.Glog(LogF, fmt.Sprintf("main exit."))
	return true
}

func (ms *MServer) Register(m *module.MetaParaMstHeartAddBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register node %v %v:%v", MstId, Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/mst/heart/add?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)

	loc, _ := time.LoadLocation("Local")
	timeLayout := "2006-01-02 15:04:05"
	stheTime, _ := time.ParseInLocation(timeLayout, m.StartTime, loc)
	sst := stheTime.Unix()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	est := etheTime.Unix()

	if est-sst <= 60 {
		m.Duration = fmt.Sprintf("%vs", est-sst)
	} else if est-sst <= 3600 {
		m.Duration = fmt.Sprintf("%vmin", (est-sst)/60)
	} else {
		m.Duration = fmt.Sprintf("%vh", (est-sst)/3600)
	}

	m.FlowNum = fmt.Sprint(len(flowspool.MemMap))
	jsonstr0, err := json.Marshal(m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		return false
	}
	jsonstr, err := util.Api_RequestPost(url, string(jsonstr0))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn1 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn1.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn1.Status_Code))
		return false
	}
	return true
}

func (ms *MServer) RegisterRemove(m *module.MetaParaMstHeartAddBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register remove node %v %v:%v", MstId, Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/mst/heart/rm?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	jsonstr0, err := json.Marshal(m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
		return false
	}
	jsonstr, err := util.Api_RequestPost(url, string(jsonstr0))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn1 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn1.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn1.Status_Code))
		return false
	}
	return true
}

type MyWork struct {
	Id       string "id"
	MstId    string "mstid"
	FlowId   string "flowid"
	Name     string "The Name of a process"
	WP       *workpool.WorkPool
	Mmw      *MetaMyWork
	StopFlag string
	JobList  []interface{}
}

func (workPool *MyWork) DoWork(workRoutine int) {
	var waitGroup util.WaitGroupWrapper
	m := new(module.MetaMstFlowRoutineHeartBean)
	u1 := uuid.Must(uuid.NewV4(),nil)
	m.Id = fmt.Sprint(u1)
	m.FlowId = workPool.Mmw.FlowId
	m.MstId = MstId
	m.RoutineId = workPool.Id
	m.Ip = Ip
	m.Port = Port
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.StartTime = fmt.Sprint(timeStr)

	workPool.JobList = make([]interface{}, 0)
	registerRoutineChannel := make(chan interface{})
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	waitGroup.Wrap(func() {
		workPool.RoutineRegister(m)
		st := time.Now().Unix()
		et := time.Now().Unix()
		for {
			select {
			case <-registerRoutineChannel:
				et = time.Now().Unix()
				if et-st > 30 {
					ret := workPool.RoutineRegister(m)
					if !ret {
						glog.Glog(LogF, fmt.Sprintf("register mst % routine %v fail.", MstId, workPool.Id))
					}
					st = time.Now().Unix()
				}
			case <-exitChan:
				glog.Glog(LogF, fmt.Sprintf("routine force stop."))
				workPool.StopFlag = "1"
				workPool.RoutineRegisterRemove(m)
				return
			default:
			}
		}
	})
	waitGroup.Wrap(func() {
		for {
			if workPool.StopFlag == "1" {
				break
			}
			registerRoutineChannel <- 2
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(60)
			time.Sleep(time.Duration(ri) * time.Second)
		}
	})
	glog.Glog(LogF, fmt.Sprintf("%s", workPool.Name))
	glog.Glog(LogF, fmt.Sprintf("*******> WR: %d  QW: %d  AR: %d", workRoutine, workPool.WP.QueuedWork(), workPool.WP.ActiveRoutines()))

	nmq := NewFlowMgr(workPool.Mmw.FlowId, workPool.Mmw.ApiServerIp, workPool.Mmw.ApiServerPort, MstId, workPool.Mmw.HomeDir, workPool.Mmw.MstIp, workPool.Mmw.MstPort, workPool.Mmw.AccessToken, workPool.Id)
	for {
		glog.Glog(LogF, "workRoutine:"+workPool.Name+",flowid:"+workPool.Mmw.FlowId+",mst:"+MstId+" working...")
		if workPool.StopFlag == "1" {
			glog.Glog(LogF, "workRoutine:"+workPool.Name+",flowid:"+workPool.Mmw.FlowId+",mst:"+MstId+" exit.")
			break
		}
		var wg util.WaitGroupWrapper
		wg.Wrap(func() { nmq.checkGo() })
		wg.Wrap(func() { nmq.checkPending() })
		wg.Wait()
		registerRoutineChannel <- 1
		rand.Seed(time.Now().UnixNano())
		ri := rand.Intn(10)
		time.Sleep(time.Duration(ri) * time.Second)
	}
	exitChan <- 1

	close(registerRoutineChannel)
	waitGroup.Wait()
	workPool.RoutineRegisterRemove(m)
}

func (workPool *MyWork) RoutineRegister(m *module.MetaMstFlowRoutineHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register %v %v routine %v:%v", MstId, workPool.Id, Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/mst/flow/routine/heart/add?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	m.Lst = workPool.JobList
	m.JobNum = fmt.Sprint(len(workPool.JobList))

	loc, _ := time.LoadLocation("Local")
	timeLayout := "2006-01-02 15:04:05"
	stheTime, _ := time.ParseInLocation(timeLayout, m.StartTime, loc)
	sst := stheTime.Unix()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	est := etheTime.Unix()

	if est-sst <= 60 {
		m.Duration = fmt.Sprintf("%vs", est-sst)
	} else if est-sst <= 3600 {
		m.Duration = fmt.Sprintf("%vmin", (est-sst)/60)
	} else {
		m.Duration = fmt.Sprintf("%vh", (est-sst)/3600)
	}

	jsonstr0, err := json.Marshal(m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	jsonstr, err := util.Api_RequestPost(url, string(jsonstr0))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn1 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn1.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn1.Status_Code))
		return false
	}
	return true
}

func (workPool *MyWork) RoutineRegisterRemove(m *module.MetaMstFlowRoutineHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register remove %v %v %v routine %v:%v", MstId, workPool.Id, Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/mst/flow/routine/heart/rm?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	jsonstr0, err := json.Marshal(m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	jsonstr, err := util.Api_RequestPost(url, string(jsonstr0))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn1 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn1.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn1.Status_Code))
		return false
	}
	return true
}
