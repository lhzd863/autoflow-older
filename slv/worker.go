package slv

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/satori/go.uuid"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/gproto"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"
)

var (
	LogF          string
	HomeDir       string
	AccessToken   string
	ApiServerIp   string
	ApiServerPort string
	jobpool       = db.NewMemDB()
	ProcessNum    int
	Ip            string
	Port          string
	WorkerId      string
)

// 业务实现方法的容器
type SServer struct {
	sync.RWMutex
	waitGroup util.WaitGroupWrapper
}

func NewSServer(paraMap map[string]interface{}) *SServer {
	WorkerId = paraMap["workerid"].(string)
	Ip = paraMap["ip"].(string)
	Port = paraMap["port"].(string)
	HomeDir = paraMap["homedir"].(string)
	AccessToken = paraMap["accesstoken"].(string)
	ApiServerIp = paraMap["apiserverip"].(string)
	ApiServerPort = paraMap["apiserverport"].(string)
	NumStr := paraMap["processnum"].(string)
	ProcessNum, _ = strconv.Atoi(NumStr)
	LogF = HomeDir + "/ms_${" + util.ENV_VAR_DATE + "}.log"
	m := &SServer{}
	return m

}

func (s *SServer) DoCmd(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	fmt.Println("MD5方法请求JSON:" + in.JsonStr)
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}

func (s *SServer) Ping(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "{}"}, nil
}

func (s *SServer) JobStart(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mjs := new(module.MetaJobWorkerBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mjs)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
	}
	err = s.Executesh(mjs)
	var status_code int32
	status_code = 200
	status_txt := ""
	if err != nil {
		status_code = 799
		status_txt = fmt.Sprint(err)
	}
	fmt.Println("MD5方法请求JSON:" + in.JsonStr)
	return &gproto.Res{Status_Txt: status_txt, Status_Code: status_code, Data: "{}"}, nil
}

func (s *SServer) Executesh(job *module.MetaJobWorkerBean) error {
	logf := fmt.Sprintf("%v/LOG/%v/%v", HomeDir, job.FlowId, job.Sys)
	m := new(module.MetaSystemWorkerRoutineJobRunningHeartBean)
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)
	m.Sys = job.Sys
	m.Job = job.Job
	m.Ip = Ip
	m.Port = Port
	m.WorkerId = WorkerId
	timeStr0 := time.Now().Format("2006-01-02 15:04:05")
	m.StartTime = timeStr0
	exitChan := make(chan int)
	StopFlag := 0
	go func() {
		<-exitChan
		StopFlag = 1
		s.WorkerJobRunningRegisterRemove(m)
	}()
	go func() {
		st := time.Now().Unix()
		et := time.Now().Unix()
		for {
			if StopFlag == 1 {
				break
			}
			et = time.Now().Unix()
			if et-st > 30 {
				s.WorkerJobRunningRegister(m)
				st = time.Now().Unix()
			}
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(2)
			time.Sleep(time.Duration(ri) * time.Second)
		}
	}()
	exist, err := util.PathExists(logf)
	if err != nil {
		exitChan <- 1
		return fmt.Errorf("failed to path exists: %v", err)
	}
	if !exist {
		os.MkdirAll(logf, os.ModePerm)
	}
	mjm := new(module.MetaJobMemBean)
	mjm.Id = job.Id
	mjm.Sys = job.Sys
	mjm.Job = job.Job
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	mjm.CreateTime = timeStr
	mjm.Enable = "1"
	jobpool.Add(job.Id, mjm)
	defer jobpool.Remove(job.Id)
	for i := 0; i < len(job.Cmd); i++ {
		timeStr := time.Now().Format("20060102150405")
		jobLogF := fmt.Sprintf("%v/%v_%v_%v_%v.log", logf, strings.ToLower(job.Job), i, job.RunningTime, timeStr)
		c := job.Cmd[i].(string)
		for j := len(job.Parameter) - 1; j >= 0; j-- {
			kv := new(module.KVBean)
			err := json.Unmarshal([]byte(job.Parameter[j].(string)), &kv)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("parse kvbean error.%v", err))
			}
			//repalce variable
			vt := "\\$\\{" + kv.K + "\\}"
			reg := regexp.MustCompile(vt)
			c = reg.ReplaceAllString(c, kv.V)
		}
		for j := 0; j < len(job.Parameter); j++ {
			kv := new(module.KVBean)
			err := json.Unmarshal([]byte(job.Parameter[j].(string)), &kv)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("parse kvbean error.%v", err))
			}
			//set env
			err = os.Setenv(kv.K, kv.V)
			glog.Glog(LogF, fmt.Sprintf("%v = %v", kv.K, kv.V))
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("job %v %v set env %v=%v error.%v", job.Sys, job.Job, kv.K, kv.V, err))
			}
		}
		cmd := exec.Command("/bin/bash", "-c", c)
		mjmt := jobpool.Get(job.Id).(*module.MetaJobMemBean)
		mjmt.Cmd = cmd
		timeStr = time.Now().Format("2006-01-02 15:04:05")
		mjmt.UpdateTime = timeStr
		jobpool.Add(job.Id, mjmt)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			_, cfile, cline, _ := runtime.Caller(1)
			glog.Glog(jobLogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
			exitChan <- 1
			return fmt.Errorf("%v", err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			_, cfile, cline, _ := runtime.Caller(1)
			glog.Glog(jobLogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
			exitChan <- 1
			return fmt.Errorf("%v", err)
		}
		cmd.Start()
		reader := bufio.NewReader(stdout)
		logarr := make([]string, 0)
		logChan := make(chan int)
		n := new(module.MetaParaFlowJobLogAddBean)
		n.Sys = job.Sys
		n.Job = job.Job
		n.FlowId = job.FlowId
		n.SServer = WorkerId
		n.Sip = Ip
		n.Sport = Port
		n.Step = fmt.Sprint(i)
		n.Cmd = c
		u2 := uuid.Must(uuid.NewV4())
		n.Id = fmt.Sprint(u2)
		n.StartTime = timeStr
		tstopflag := 0
		var wg util.WaitGroupWrapper
		go func() {
			<-logChan
			tstopflag = 1
		}()
		wg.Wrap(func() {
			st := time.Now().Unix()
			et := time.Now().Unix()
			n.Content = make([]string, 0)
			s.JobLogAppend(n)
			for {
				if tstopflag == 1 {
					break
				}
				et = time.Now().Unix()
				if len(logarr) >= 100 || (et-st >= 10 && len(logarr) > 0) {
					s.Lock()
					tarr := logarr
					logarr = make([]string, 0)
					s.Unlock()
					n.Content = tarr
					s.JobLogAppend(n)
					st = time.Now().Unix()
				}
				rand.Seed(time.Now().UnixNano())
				ri := rand.Intn(2)
				time.Sleep(time.Duration(ri) * time.Second)
			}
		})
		go func() {
			for {
				line, err2 := reader.ReadString('\n')
				if err2 != nil || io.EOF == err2 {
					break
				}
				glog.Glog(jobLogF, fmt.Sprintf("%v", line))
				logarr = append(logarr, line)
			}
		}()
		readererr := bufio.NewReader(stderr)
		go func() {
			for {
				line, err2 := readererr.ReadString('\n')
				if err2 != nil || io.EOF == err2 {
					break
				}
				glog.Glog(jobLogF, fmt.Sprintf("%v", line))
				logarr = append(logarr, line)
			}
		}()

		cmd.Wait()
		logChan <- 1
		retcd := string(fmt.Sprintf("%v", cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()))
		retcd = strings.Replace(retcd, " ", "", -1)
		retcd = strings.Replace(retcd, "\n", "", -1)
		timeStr1 := time.Now().Format("2006-01-02 15:04:05")
		n.EndTime = timeStr1
		wg.Wait()
		n.Content = logarr
		if retcd != "0" {
			n.ExitCode = retcd
			s.JobLogAppend(n)
			glog.Glog(jobLogF, fmt.Sprintf("%v", retcd))
			exitChan <- 1
			return fmt.Errorf("%v", retcd)
		}
		n.ExitCode = "0"
		s.JobLogAppend(n)
	}
	exitChan <- 1
	return nil
}

func (s *SServer) JobLogAppend(n *module.MetaParaFlowJobLogAddBean) bool {
	glog.Glog(LogF, fmt.Sprintf("node %v, %v", Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/flow/job/log/append?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)

	jsonstr0, err := json.Marshal(n)
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

func (s *SServer) WorkerJobRunningRegister(m *module.MetaSystemWorkerRoutineJobRunningHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register node %v, %v", Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/routine/job/running/heart/add?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
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

func (s *SServer) WorkerJobRunningRegisterRemove(m *module.MetaSystemWorkerRoutineJobRunningHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register node %v, %v", Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/routine/job/running/heart/rm?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
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

// 为server定义 DoMD5 方法 内部处理请求并返回结果
// 参数 (context.Context[固定], *test.Req[相应接口定义的请求参数])
// 返回 (*test.Res[相应接口定义的返回参数，必须用指针], error)
func (s *SServer) JobStop(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mjs := new(module.MetaJobWorkerBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mjs)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
	}
	for k := range jobpool.MemMap {
		jpo := (jobpool.MemMap[k]).(*module.MetaJobMemBean)
		if k == mjs.Id {
			jpo.Cmd.Process.Kill()
		}
	}
	glog.Glog(LogF, fmt.Sprintf("%v,%v,%v has stopped.", mjs.Id, mjs.Sys, mjs.Job))
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}

func (s *SServer) JobStatus(ctx context.Context, in *gproto.Req) (*gproto.Res, error) {
	mjs := new(module.MetaJobWorkerBean)
	err := json.Unmarshal([]byte(in.JsonStr), &mjs)
	if err != nil {
		_, cfile, cline, _ := runtime.Caller(1)
		glog.Glog(LogF, fmt.Sprintf("%v %v %v", cfile, cline, err))
	}
	var jsonstr []byte
	for k := range jobpool.MemMap {
		jpo := (jobpool.MemMap[k]).(*module.MetaJobMemBean)
		if k == mjs.Id {
			jsonstr, _ = json.Marshal(jpo)
		}
	}
	return &gproto.Res{Status_Txt: "", Status_Code: 200, Data: string(jsonstr)}, nil
}

func (s *SServer) Main() bool {
	var wg util.WaitGroupWrapper
	StopFlag := 0
	m := new(module.MetaWorkerHeartBean)
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)
	timeStr0 := time.Now().Format("2006-01-02 15:04:05")
	m.StartTime = timeStr0
	m.WorkerId = WorkerId
	m.Ip = Ip
	m.Port = Port
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		StopFlag = 1
		s.RegisterRemove(m)
		exitChan <- 1
	}()
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		st := time.Now().Unix()
		et := time.Now().Unix()
		for {
			if StopFlag == 1 {
				return
			}
			et = time.Now().Unix()
			if et-st > 30 {
				ret := s.Register(m)
				if !ret {
					glog.Glog(LogF, "register worker fail.")
				}
				st = time.Now().Unix()
			}
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(2)
			time.Sleep(time.Duration(ri) * time.Second)
		}
	}()

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
	gproto.RegisterSlaverServer(ss, &SServer{})
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
	ss.Stop()
	wg.Wait()
	return true
}

func (s *SServer) Register(m *module.MetaWorkerHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register node %v, %v", Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/heart/add?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
	m.MaxCnt = fmt.Sprint(ProcessNum)
	m.RunningCnt = fmt.Sprint(len(jobpool.MemMap))
	m.CurrentExecCnt = "0"
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

func (s *SServer) RegisterRemove(m *module.MetaWorkerHeartBean) bool {
	glog.Glog(LogF, fmt.Sprintf("Register node %v, %v", Ip, Port))
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/heart/rm?accesstoken=%v", ApiServerIp, ApiServerPort, AccessToken)
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
