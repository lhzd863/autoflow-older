package leader

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/lhzd863/autoflow/glog"
	"github.com/lhzd863/autoflow/gproto"
	"github.com/lhzd863/autoflow/module"
	"github.com/lhzd863/autoflow/util"
	uuid "github.com/satori/go.uuid"
)

type FlowMgr struct {
	LeaderId                     string
	Lip                          string
	Lport                        string
	FlowId                       string
	RoutineId                    string
	ApiServerIp                  string
	ApiServerPort                string
	Name                         string
	Timestamp                    int64
	Attempts                     uint16
	StopFlag                     bool
	MaxJobCount                  int64
	QueueDir                     string
	SleepTime                    int64
	LogF                         string
	HomeDir                      string
	AccessToken                  string
	goStatusChannel              chan []interface{}
	shutdownGoStatusChannel      chan interface{}
	pendingStatusChannel         chan []interface{}
	shutdownPendingStatusChannel chan interface{}
	sync.RWMutex
}

func NewFlowMgr(flowid string, apiserverip string, apiserverport string, LeaderId string, homeDir string, leaderip string, leaderport string, accesstoken string, routineid string) *FlowMgr {
	return &FlowMgr{
		LeaderId:      LeaderId,
		Lip:           leaderip,
		Lport:         leaderport,
		FlowId:        flowid,
		RoutineId:     routineid,
		ApiServerIp:   apiserverip,
		ApiServerPort: apiserverport,
		Name:          "leader",
		Timestamp:     time.Now().Unix(),
		Attempts:      0,
		StopFlag:      false,
		QueueDir:      homeDir + "/queue",
		SleepTime:     10,
		MaxJobCount:   100,
		LogF:          homeDir + "/leader_" + flowid + ".log",
		HomeDir:       homeDir,
		AccessToken:   accesstoken,
	}
}

//check status go
func (m *FlowMgr) checkGo() bool {
	glog.Glog(m.LogF, "Checking Status Go.")
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/get/go?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	u1 := uuid.Must(uuid.NewV4(), nil)
	para := fmt.Sprintf("{\"id\":\"%v\",\"flowid\":\"%v\",\"ishash\":\"0\",\"status\":\"%v\"}", u1, m.FlowId, util.STATUS_AUTO_GO)

	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	if retbn.Status_Code != 200 {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v post url return status code:%v", cfile, cline, retbn.Status_Code))
		return false
	}
	if retbn.Status_Txt == "-1" {
		glog.Glog(LogF, fmt.Sprintf("no ring id use."))
		return false
	}
	if retbn.Data == nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v get pending status job err.", cfile, cline))
		return false
	}
	retarr := (retbn.Data).([]interface{})
	for i := 0; i < len(retarr); i++ {
		v := retarr[i].(map[string]interface{})
		if v["enable"] != "1" {
			glog.Glog(m.LogF, fmt.Sprintf("Job %v %v is not enabled,wait for next time.", v["job"], v["enable"]))
			continue
		}
		//execute job

		var waitGroup util.WaitGroupWrapper
		jobv := m.jobInfo(v["sys"].(string), v["job"].(string))
		if len(jobv) > 0 {
			if v["jobtype"].(string) == "V" {
				//var waitGroup util.WaitGroupWrapper
				waitGroup.Wrap(func() { m.invokeVirtualJob(jobv[0].(map[string]interface{})) })
			} else {
				//var waitGroup util.WaitGroupWrapper
				err = m.workerExecApplication(v["wserver"].(string))
				if err != nil {
					glog.Glog(m.LogF, fmt.Sprint(err))
					continue
				}
				waitGroup.Wrap(func() { m.invokeRealJob(jobv[0].(map[string]interface{})) })
			}
		}
	}
	if len(retarr) == 0 {
		glog.Glog(m.LogF, fmt.Sprintf("no go job.%v", retbn.Status_Txt))
	} else {
		m.goRemoveRing(fmt.Sprint(u1))
	}
	return true
}

func (m *FlowMgr) goRemoveRing(id string) bool {
	glog.Glog(m.LogF, "go Remove ring.")
	url := fmt.Sprintf("http://%v:%v/api/v1/system/ring/go/rm?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"id\":\"%v\"}", id)

	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	if retbn.Status_Code != 200 {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v post url return status code:%v", cfile, cline, retbn.Status_Code))
		return false
	}
	return true
}

func (m *FlowMgr) invokeVirtualJob(job map[string]interface{}) {
	glog.Glog(m.LogF, fmt.Sprintf("Virtual %v %v", job["sys"], job["job"]))
	//change job status
	//u1 := uuid.Must(uuid.NewV4())
	//retes := m.jobStatusUpdate(job,util.SYS_STATUS_DONE,"et",fmt.Sprint(u1))
	//stream job
	//if retes {
	//	_ = m.streamJob(job)
	//}
}

func (m *FlowMgr) jobStepCmd(sys string, job string) []interface{} {
	retarr := make([]interface{}, 0)
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/cmd/getall?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, sys, job)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return err.%v", retbn.Status_Txt))
		return retarr
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprintf("get pending status job err."))
		return retarr
	}
	glog.Glog(LogF, fmt.Sprint(retbn.Data))
	return (retbn.Data).([]interface{})
}

func (m *FlowMgr) jobParameter(sys string, job string) []interface{} {
	retarr := make([]interface{}, 0)
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/parameter/getall?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, sys, job)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return err.%v", retbn.Status_Txt))
		return retarr
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprintf("get pending status job err."))
		return retarr
	}
	return retbn.Data.([]interface{})
}

func (m *FlowMgr) streamJob(job map[string]interface{}) bool {
	glog.Glog(m.LogF, fmt.Sprintf("%v.%v stream job.", job["sys"], job["job"]))
	m.Lock()
	defer m.Unlock()
	streamarr := m.jobStream(job["sys"].(string), job["job"].(string))
	if len(streamarr) > 0 {
		for _, v := range streamarr {
			v1 := v.(map[string]interface{})
			if v1["enable"] != "1" {
				glog.Glog(m.LogF, fmt.Sprintf("%v.%v stream %v.%v enable %v is not enabled,wait for next time.", v1["streamsys"], v1["streamjob"], v1["sys"], v1["job"], v1["enable"]))
				continue
			}
			glog.Glog(m.LogF, fmt.Sprintf("%v.%v stream %v.%v.", job["sys"], job["job"], v1["sys"], v1["job"]))
			//fail retry 3 times
			for j := 0; j < 3; j++ {
				err := m.jobStreamJob(v1["sys"].(string), v1["job"].(string))
				if err == nil {
					glog.Glog(m.LogF, fmt.Sprintf("%v.%v stream %v.%v successfully.", v1["streamsys"], v1["streamjob"], v1["sys"], v1["job"]))
					break
				} else {
					glog.Glog(m.LogF, fmt.Sprintf("%v.%v stream %v.%v fail %v times,%v.", v1["streamsys"], v1["streamjob"], v1["sys"], v1["job"], j, err))
				}
				rand.Seed(time.Now().UnixNano())
				ri := rand.Intn(10)
				time.Sleep(time.Duration(ri) * time.Millisecond)
			}
		}
	}

	return true
}

func (m *FlowMgr) jobStream(sys string, job string) []interface{} {
	retarr := make([]interface{}, 0)
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/stream/job/get?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, sys, job)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return err.%v", retbn.Status_Txt))
		return retarr
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprintf("get pending status job err."))
		return retarr
	}
	return retbn.Data.([]interface{})
}

func (m *FlowMgr) jobStreamJob(sys string, job string) error {
	glog.Glog(LogF, fmt.Sprintf("%v.%v stream job.", sys, job))
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/stream/job?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, sys, job)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return err.%v", retbn.Status_Txt))
		return errors.New(retbn.Status_Txt)
	}
	return nil
}

func (m *FlowMgr) invokeRealJob(job map[string]interface{}) {
	glog.Glog(m.LogF, fmt.Sprintf("exec %v.%v on slave %v [%v:%v].", job["sys"], job["job"], job["wserver"], job["wip"], job["wport"]))
	SFlag := 0
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	var waitGroup util.WaitGroupWrapper
	exitChan := make(chan int)
	mf := new(module.MetaParaSystemLeaderFlowRoutineJobRunningHeartAddBean)
	u1 := uuid.Must(uuid.NewV4(), nil)
	mf.Id = fmt.Sprint(u1)
	_, ok := job["sys"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job sys att is null"))
		return
	}
	mf.Sys = job["sys"].(string)

	_, ok = job["job"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job job att is null"))
		return
	}
	mf.Job = job["job"].(string)

	_, ok = job["wip"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job wip att is null"))
		return
	}
	mf.Wip = job["wip"].(string)

	_, ok = job["wport"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job wport att is null"))
		return
	}
	mf.Wport = job["wport"].(string)

	_, ok = job["wserver"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job wserver att is null"))
		return
	}
	mf.WorkerId = job["wserver"].(string)
	mf.StartTime = timeStr

	waitGroup.Wrap(func() {
		<-exitChan
		SFlag = 1
		m.RegisterRemove(mf)
	})

	waitGroup.Wrap(func() {
		st := time.Now().Unix()
		et := time.Now().Unix()
		for {
			if SFlag == 1 {
				break
			}
			et = time.Now().Unix()
			if et-st > 30 {
				ret := m.Register(mf)
				if !ret {
					glog.Glog(LogF, "register leader fail.")
				}
				st = time.Now().Unix()
			}
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(2)
			time.Sleep(time.Duration(ri) * time.Second)
		}
	})

	mjwb := new(module.MetaJobWorkerBean)
	u1 = uuid.Must(uuid.NewV4(), nil)
	mjwb.Id = fmt.Sprint(u1)
	mjwb.FlowId = m.FlowId
	mjwb.Sys = job["sys"].(string)
	mjwb.Job = job["job"].(string)

	_, ok = job["retry"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job retry att is null"))
		return
	}
	retry, err := strconv.Atoi(job["retry"].(string))
	if err != nil {
		glog.Glog(m.LogF, fmt.Sprintf("job retry is nil ,will set default 1."))
		retry = 1
	}
	mjwb.Retry = retry

	_, ok = job["alert"]
	if !ok {
		glog.Glog(m.LogF, fmt.Sprintf("job alert att is null"))
		return
	}
	mjwb.Alert = job["alert"].(string)
	mjwb.Status = util.STATUS_AUTO_RUNNING
	mjwb.StartTime = timeStr
	mjwb.WorkerIp = job["wip"].(string)
	mjwb.WorkerPort = job["wport"].(string)
	arrstep := m.jobStepCmd(job["sys"].(string), job["job"].(string))
	mjwb.Cmd = make([]interface{}, 0)
	for i := 0; i < len(arrstep); i++ {
		ast := arrstep[i].(map[string]interface{})

		_, ok = ast["cmd"]
		if !ok {
			glog.Glog(m.LogF, fmt.Sprintf("job step cmd att is null"))
			return
		}
		mjwb.Cmd = append(mjwb.Cmd, ast["cmd"].(string))
	}
	arrparameter := m.jobParameter(job["sys"].(string), job["job"].(string))
	mjwb.Parameter = make([]interface{}, 0)
	for i := 0; i < len(arrparameter); i++ {
		ast := arrparameter[i].(map[string]interface{})
		b := new(module.KVBean)
		_, ok = ast["key"]
		if !ok {
			glog.Glog(m.LogF, fmt.Sprintf("key att is null"))
			return
		}
		b.K = ast["key"].(string)

		_, ok = ast["val"]
		if !ok {
			glog.Glog(m.LogF, fmt.Sprintf("val att is null"))
			return
		}
		b.V = ast["val"].(string)
		jsonstr0, err := json.Marshal(b)
		if err != nil {
			glog.Glog(m.LogF, fmt.Sprint(err))
			continue
		}
		mjwb.Parameter = append(mjwb.Parameter, string(jsonstr0))
	}
	jsonstr, err := json.Marshal(mjwb)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	//m.jobStatusUpdate(job,util.STATUS_AUTO_RUNNING,"st",fmt.Sprint(u1))

	err = m.updateStatusRunning(job["sys"].(string), job["job"].(string))
	if err != nil {
		glog.Glog(m.LogF, fmt.Sprintf("could not update status: %v", err))
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	// 建立连接到gRPC服务
	conn, err := grpc.Dial(job["wip"].(string)+":"+job["wport"].(string), grpc.WithInsecure())
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v did not connect: %v", cfile, cline, err))
		err = m.updateStatusEnd(job["sys"].(string), job["job"].(string), util.STATUS_AUTO_FAIL)
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := gproto.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 调用gRPC接口
	tr, err := t.JobStart(ctx, &gproto.Req{JsonStr: string(jsonstr)})
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v could not greet: %v", cfile, cline, err))
		err = m.updateStatusEnd(job["sys"].(string), job["job"].(string), util.STATUS_AUTO_FAIL)
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	//change job status
	if tr.Status_Code == 200 {
		err = m.updateStatusEnd(job["sys"].(string), job["job"].(string), util.STATUS_AUTO_SUCC)
	} else {
		glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
		err = m.updateStatusEnd(job["sys"].(string), job["job"].(string), util.STATUS_AUTO_FAIL)
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	//stream job
	if err != nil {
		glog.Glog(m.LogF, fmt.Sprint(err))
		exitChan <- 1
		waitGroup.Wait()
		return
	}
	err = m.workerExecApplicationLogout(job["wserver"].(string))
	if err != nil {
		glog.Glog(m.LogF, fmt.Sprint(err))
	}
	exitChan <- 1
	waitGroup.Wait()
	m.streamJob(job)
}

func (m *FlowMgr) Register(mf *module.MetaParaSystemLeaderFlowRoutineJobRunningHeartAddBean) bool {
	// glog.Glog(LogF, fmt.Sprintf("Register job %v %v:%v", LeaderId, Ip, Port))
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	url := fmt.Sprintf("http://%v:%v/api/v1/leader/flow/routine/job/running/heart/add?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	mf.LeaderId = m.LeaderId
	mf.FlowId = m.FlowId
	mf.RoutineId = m.RoutineId
	mf.Mip = m.Lip
	mf.Mport = m.Lport

	loc, _ := time.LoadLocation("Local")
	timeLayout := "2006-01-02 15:04:05"
	stheTime, _ := time.ParseInLocation(timeLayout, mf.StartTime, loc)
	sst := stheTime.Unix()
	etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	est := etheTime.Unix()

	if est-sst <= 60 {
		mf.Duration = fmt.Sprintf("%vs", est-sst)
	} else if est-sst <= 3600 {
		mf.Duration = fmt.Sprintf("%vmin", (est-sst)/60)
	} else {
		mf.Duration = fmt.Sprintf("%vh", (est-sst)/3600)
	}
	jsonstr0, err := json.Marshal(mf)
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

func (m *FlowMgr) RegisterRemove(mf *module.MetaParaSystemLeaderFlowRoutineJobRunningHeartAddBean) bool {
	// glog.Glog(LogF, fmt.Sprintf("Register job %v %v:%v", LeaderId, Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/leader/flow/routine/job/running/heart/rm?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	mf.LeaderId = m.LeaderId
	mf.FlowId = m.FlowId
	mf.RoutineId = m.RoutineId

	jsonstr0, err := json.Marshal(mf)
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

func (m *FlowMgr) updateStatusRunning(sys string, job string) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/update/start?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\",\"status\":\"%v\"}", m.FlowId, sys, job, util.STATUS_AUTO_RUNNING)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return errors.New(fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
	}
	return nil
}

func (m *FlowMgr) updateStatusEnd(sys string, job string, status string) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/update/end?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\",\"status\":\"%v\"}", m.FlowId, sys, job, status)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return errors.New(fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
	}
	return nil
}

func (m *FlowMgr) updateStatus2Server(sys string, job string, status string, server string) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/update/2server?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	s := new(module.MetaParaFlowJobStatus2ServerBean)
	s.FlowId = m.FlowId
	s.Sys = sys
	s.Job = job
	s.Status = status
	s.Server = server
	jsonstr0, err := json.Marshal(s)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return err
	}
	jsonstr, err := util.Api_RequestPost(url, string(jsonstr0))
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return err
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return errors.New(fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
	}
	return nil
}

func (m *FlowMgr) workerExecApplication(workerid string) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/exec/add?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"workerid\":\"%v\"}", workerid)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(fmt.Sprintf("post url return err:%v", retbn.Status_Txt))
	}
	if retbn.Data == nil {
		return errors.New("application exec err.")
	}
	if len((retbn.Data).([]interface{})) == 0 {
		return errors.New("application 0 exec.")
	}
	return nil
}

func (m *FlowMgr) workerExecApplicationLogout(workerid string) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/exec/sub?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"workerid\":\"%v\"}", workerid)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(fmt.Sprintf("post url return err:%v", retbn.Status_Txt))
	}
	return nil
}

func (m *FlowMgr) jobInfo(sys string, job string) []interface{} {
	retarr := make([]interface{}, 0)
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/get?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, sys, job)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return retarr
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return retarr
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprintf("get pending status job err."))
		return retarr
	}
	return (retbn.Data).([]interface{})
}

//Check the status pending
func (m *FlowMgr) checkPending() bool {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	glog.Glog(m.LogF, fmt.Sprintf("%v %v %v Checking %v Status Pending.", m.LeaderId, m.FlowId, m.RoutineId, timeStr))
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/get/pending?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	u1 := uuid.Must(uuid.NewV4(), nil)
	para := fmt.Sprintf("{\"id\":\"%v\",\"flowid\":\"%v\",\"ishash\":\"0\",\"status\":\"%v\"}", u1, m.FlowId, util.STATUS_AUTO_PENDING)
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return false
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprint("get pending status job err."))
		return false
	}
	retarr := (retbn.Data).([]interface{})
	for i := 0; i < len(retarr); i++ {
		v := retarr[i].(map[string]interface{})
		if v["enable"] != "1" {
			glog.Glog(m.LogF, fmt.Sprintf("Job %v %v is not enabled,wait for next time.", v["job"], v["enable"]))
			continue
		}
		if err = m.isDependencyOk(v); err != nil {
			glog.Glog(m.LogF, fmt.Sprintf("Dependency is not ok,%v", err))
			continue
		}
		if err = m.isTimeWindowOk(v); err != nil {
			glog.Glog(m.LogF, fmt.Sprintf("Timewindow is not ok,%v", err))
			continue
		}
		if err = m.isCmdOk(v); err != nil {
			glog.Glog(m.LogF, fmt.Sprintf("Cmd is not ok,%v", err))
			continue
		}
		if err = m.submitJob(v); err != nil {
			glog.Glog(m.LogF, fmt.Sprint("Submit job is not ok.", err))
		}
	}
	if len(retarr) == 0 {
		glog.Glog(m.LogF, fmt.Sprintf("no pending job.%v", retbn.Status_Txt))
	} else {
		m.pendingRemoveRing(fmt.Sprint(u1))
	}
	return true
}

func (m *FlowMgr) pendingRemoveRing(id string) bool {
	glog.Glog(m.LogF, "pending remove ring.")
	url := fmt.Sprintf("http://%v:%v/api/v1/system/ring/pending/rm?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"id\":\"%v\"}", id)

	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(m.LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
		return false
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return false
	}
	return true
}

func (m *FlowMgr) submitJob(job map[string]interface{}) error {
	url0 := fmt.Sprintf("http://%v:%v/api/v1/job/pool/add?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para0 := fmt.Sprintf("{\"sys\":\"%v\",\"job\":\"%v\",\"flowid\":\"%v\",\"priority\":\"%v\",\"dynamicserver\":\"%v\",\"server\":\"%v\"}", job["sys"], job["job"], m.FlowId, job["priority"], job["dynamicserver"], job["wserver"])
	jsonstr0, err := util.Api_RequestPost(url0, para0)
	if err != nil {
		return err
	}
	retbn0 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr0), &retbn0)
	if err != nil {
		return err
	}
	if retbn0.Status_Code != 200 {
		return errors.New(retbn0.Status_Txt)
	}
	glog.Glog(m.LogF, fmt.Sprintf("update job status %v.%v submit", job["sys"], job["job"]))
	url1 := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/update/submit?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para1 := fmt.Sprintf("{\"sys\":\"%v\",\"job\":\"%v\",\"flowid\":\"%v\",\"status\":\"%v\"}", job["sys"], job["job"], m.FlowId, util.STATUS_AUTO_SUBMIT)
	glog.Glog(LogF, fmt.Sprint(url1))
	jsonstr, err := util.Api_RequestPost(url1, para1)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(retbn.Status_Txt)
	}
	return nil
}

func (m *FlowMgr) processControlFile(ctlinfo interface{}) bool {
	return true
}

//Sort the control file by prior
func (m *FlowMgr) sortControlFile(fmap map[string]interface{}) []interface{} {
	var i = 0
	alen := len(fmap)
	var arr = make([]interface{}, alen)
	for _, v := range fmap {
		arr[i] = v
		i++
	}
	return arr
}

func (m *FlowMgr) getCurrentJobCount() int64 {
	return 1
}

func (m *FlowMgr) isControlFile(filename string) bool {
	arr := strings.Split(filename, ".")
	if len(arr) == 3 {
		return true
	} else {
		return false
	}
}

func (m *FlowMgr) LeaderCtlMarshal() ([]byte, error) {
	ctl := new(MetaJobCTL)
	data, err := json.Marshal(ctl)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *FlowMgr) leaderCtlUnmarshal(data []byte) (MetaJobCTL, error) {
	var ctl MetaJobCTL
	err := json.Unmarshal(data, &ctl)
	if err != nil {
		return ctl, err
	}
	return ctl, nil
}

func (m *FlowMgr) leaderCtlRead(f string) ([]byte, error) {
	fp, err := os.OpenFile(f, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	data := make([]byte, 1024)
	n, err := fp.Read(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *FlowMgr) leaderCtlWrite(f string, data []byte) error {
	fp, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (m *FlowMgr) isDependencyOk(job map[string]interface{}) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/dependency?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, job["sys"], job["job"])
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(retbn.Status_Txt)
	}
	if retbn.Data == nil {
		return errors.New("return data null.")
	}
	retarr := (retbn.Data).([]interface{})
	if len(retarr) > 0 {
		v := retarr[0].(map[string]interface{})
		return errors.New(fmt.Sprintf("%v.%v dependant %v.%v not finished, wait for next time!", v["sys"], v["job"], v["dependencysys"], v["dependencyjob"]))
	}
	return nil
}

func (m *FlowMgr) isTimeWindowOk(job map[string]interface{}) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/timewindow?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, job["sys"], job["job"])
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(retbn.Status_Txt)
	}
	if retbn.Data == nil {
		return errors.New("return data null.")
	}
	retarr := (retbn.Data).([]interface{})
	if len(retarr) > 0 {
		v := retarr[0].(map[string]interface{})
		timeHour := time.Now().Hour()
		return errors.New(fmt.Sprintf("The current hour %v does not match %v-%v the job %v %v time window, wait for next time.", timeHour, v["starthour"], v["endhour"], v["sys"], v["job"]))
	}
	return nil
}

func (m *FlowMgr) isCmdOk(job map[string]interface{}) error {
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/cmd/getall?accesstoken=%v", m.ApiServerIp, m.ApiServerPort, m.AccessToken)
	para := fmt.Sprintf("{\"flowid\":\"%v\",\"sys\":\"%v\",\"job\":\"%v\"}", m.FlowId, job["sys"], job["job"])
	jsonstr, err := util.Api_RequestPost(url, para)
	if err != nil {
		return err
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		return err
	}
	if retbn.Status_Code != 200 {
		return errors.New(retbn.Status_Txt)
	}
	if retbn.Data == nil {
		return errors.New("return data null.")
	}
	retarr := (retbn.Data).([]interface{})
	if len(retarr) > 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("%v.%v is not exists cmd, wait for next time.", job["sys"], job["job"]))
}
