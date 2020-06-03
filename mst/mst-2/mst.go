package mst

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/jwt"
	"github.com/lhzd863/autoflow/internal/util"
	"github.com/lhzd863/autoflow/internal/workpool"
)

var (
	FlowId        string
	ApiServerIp   string
	ApiServerPort string
	LogF          string
	Msgch         chan string
	CHashRing     *Consistent
	HomeDir       string
)

type Mst struct {
	Port      string
	JwtKey    string
	FlowId    string
	waitGroup util.WaitGroupWrapper
	HttpSer   *http.Server
}

func NewMst(paraMap map[string]interface{}) *Mst {
	FlowId = paraMap["flowid"].(string)
	ApiServerIp = paraMap["apiserverip"].(string)
	ApiServerPort = paraMap["apiserverport"].(string)
	HomeDir = paraMap["homedir"].(string)
	LogF = HomeDir + "/" + FlowId + "/mgr_${" + util.ENV_VAR_DATE + "}.log"
	m := &Mst{
		Port:   paraMap["port"].(string),
		JwtKey: paraMap["jwtkey"].(string),
		FlowId: FlowId,
	}
	return m
}

func (m *Mst) StartServer() *http.Server {
	if ok, _ := util.PathExists(HomeDir + "/" + FlowId + "/LOG"); !ok {
		os.Mkdir(HomeDir+"/"+FlowId+"/LOG", os.ModePerm)
	}
	var ser *http.Server
        ser = m.HttpServer()
	m.HttpSer = ser
	m.Main()
	return ser
}

func (m *Mst) Main() {
	Msgch = make(chan string)
	CHashRing = NewConsistent()
	for i := 0; i < 3; i++ {
		si := fmt.Sprintf("%d", i)
		CHashRing.Add(NewNode(i, "m"+si, 1))
	}
	workPool := workpool.New(3, 1000)
	for j := 0; j < 3; j++ {
		mw := new(MetaMyWork)
		mw.HomeDir = HomeDir
		mw.FlowId = FlowId
		mw.ApiServerIp = ApiServerIp
		mw.ApiServerPort = ApiServerPort
		work := &MyWork{
			Id:   fmt.Sprint(j),
			Name: "A-" + fmt.Sprint(j),
			WP:   workPool,
			Mmw:  mw,
		}
		err := workPool.PostWork("name_routine", work)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("ERROR: %s\n", err))
			continue
		}
	}
}

func (m *Mst) HttpServer() *http.Server {
	glog.Glog(LogF, fmt.Sprintf("Listener port :%v", m.Port))
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/mst/job/list", JobListHandler)
	mux.HandleFunc("/mst/jobqueue/add", JobQueueAddHandler)
	mux.HandleFunc("/mst/jobqueue/get", JobQueueGetHandler)
	mux.HandleFunc("/mst/jobqueue/ls", JobQueueListHandler)
	mux.HandleFunc("/mst/jobqueue/rm", JobQueueRemoveHandler)
	server := &http.Server{
		Addr:         ":" + m.Port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      mux,
	}
	// these timeouts are absolute per server connection NOT per request
	// this means that a single persistent connection will only last N seconds
	m.waitGroup.Wrap(func() {
		err := server.ListenAndServe()
		// theres no direct way to detect this error because it is not exposed
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			glog.Glog(LogF, fmt.Sprintf("ERROR: http.ListenAndServe() - %s", err.Error()))
		}
	})
	return server
}

func pingHandler(response http.ResponseWriter, request *http.Request) {
	err := globalOauth(response, request)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint("Not Authorized.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("Not Authorized.%v", err), nil)
		return
	}
	util.ApiResponse(response, 200, "", nil)
}

func JobListHandler(response http.ResponseWriter, request *http.Request) {
	util.ApiResponse(response, 200, "", nil)
}

func JobQueueAddHandler(response http.ResponseWriter, request *http.Request) {
	glog.Glog(LogF, "Method:"+fmt.Sprint(request.Method))
	if request.Method != "POST" {
		glog.Glog(LogF, "706 method is not post")
		util.ApiResponse(response, 706, "method is not post", nil)
		return
	}
	err := globalOauth(response, request)
	if err != nil {
		glog.Glog(LogF, "700 Not Authorized")
		util.ApiResponse(response, 700, "Not Authorized", nil)
		return
	}
	reqParams, err := util.NewReqParams(request)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Failed to decode: %v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("Failed to decode.%v", err), nil)
		return
	}
	sys, err := reqParams.Get("sys")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get sys parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get sys parameter failed.%v", err), nil)
		return
	}
	job, err := reqParams.Get("job")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get job parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get job parameter failed.%v", err), nil)
		return
	}
	txts, err := reqParams.Get("txts")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get txts parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get txts parameter failed.%v", err), nil)
		return
	}

	queuebean := new(MetaJobQueue)
	queuebean.Sys = sys
	queuebean.Job = job
	queuebean.Txts = txts
	bt := db.NewBoltDB(HomeDir+"/"+FlowId+"/"+FlowId+".db", util.SYS_TB_JOB_QUEUE)
	defer bt.Close()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	queuebean.Cts = timeStr
	queuebean.Enable = "1"
	jsonstr, _ := json.Marshal(queuebean)
	err = bt.Set(queuebean.Sys+"."+queuebean.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	util.ApiResponse(response, 200, "", nil)
}

func JobQueueGetHandler(response http.ResponseWriter, request *http.Request) {
	glog.Glog(LogF, "Method:"+fmt.Sprint(request.Method))
	if request.Method != "POST" {
		glog.Glog(LogF, "706 method is not post")
		util.ApiResponse(response, 706, "method is not post", nil)
		return
	}
	err := globalOauth(response, request)
	if err != nil {
		glog.Glog(LogF, "700 Not Authorized")
		util.ApiResponse(response, 700, "Not Authorized", nil)
		return
	}
	reqParams, err := util.NewReqParams(request)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Failed to decode: %v", err))
		return
	}
	sys, err := reqParams.Get("sys")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get sys parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get sys parameter failed.%v", err), nil)
		return
	}
	job, err := reqParams.Get("job")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get job parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get job parameter failed.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(HomeDir+"/"+FlowId+"/"+FlowId+".db", util.SYS_TB_JOB_QUEUE)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	qb := bt.Get(sys + "." + job)
	if qb != nil {
		jqb := new(MetaJobQueue)
		err := json.Unmarshal([]byte(qb.(string)), &jqb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, jqb)
	}
	util.ApiResponse(response, 200, "", retlst)
}

func JobQueueRemoveHandler(response http.ResponseWriter, request *http.Request) {
	glog.Glog(LogF, "Method:"+fmt.Sprint(request.Method))
	if request.Method != "POST" {
		glog.Glog(LogF, "706 method is not post")
		util.ApiResponse(response, 706, "method is not post", nil)
		return
	}
	err := globalOauth(response, request)
	if err != nil {
		glog.Glog(LogF, "700 Not Authorized")
		util.ApiResponse(response, 700, "Not Authorized", nil)
		return
	}
	reqParams, err := util.NewReqParams(request)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Failed to decode: %v", err))
		return
	}
	sys, err := reqParams.Get("sys")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get sys parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get sys parameter failed.%v", err), nil)
		return
	}
	job, err := reqParams.Get("job")
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get job parameter failed.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("get job parameter failed.%v", err), nil)
		return
	}

	bt := db.NewBoltDB(HomeDir+"/"+FlowId+"/"+FlowId+".db", util.SYS_TB_JOB_QUEUE)
	defer bt.Close()

	err = bt.Remove(sys + "." + job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("rm data err.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("rm data err.%v", err), nil)
		return
	}
	util.ApiResponse(response, 200, "", nil)
}

func JobQueueListHandler(response http.ResponseWriter, request *http.Request) {
	glog.Glog(LogF, "Method:"+fmt.Sprint(request.Method))
	if request.Method != "POST" {
		glog.Glog(LogF, "706 method is not post")
		util.ApiResponse(response, 706, "method is not post", nil)
		return
	}
	err := globalOauth(response, request)
	if err != nil {
		glog.Glog(LogF, "700 Not Authorized")
		util.ApiResponse(response, 700, "Not Authorized", nil)
		return
	}

	bt := db.NewBoltDB(HomeDir+"/"+FlowId+"/"+FlowId+".db", util.SYS_TB_JOB_QUEUE)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jqb := new(MetaJobQueue)
			err := json.Unmarshal([]byte(v1.(string)), &jqb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, jqb)
		}
	}
	util.ApiResponse(response, 200, "", retlst)
}

// Global Filter
func globalOauth(response http.ResponseWriter, request *http.Request) error {
	reqParams, err := util.NewReqParams(request)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Failed to decode: %v", err))
		return err
	}
	tokenstring, err := reqParams.Get("accesstoken")
	if err != nil {
		glog.Glog(LogF, "no accesstoken")
		return errors.New("no accesstoken")
	}

	var claimsDecoded map[string]interface{}
	decodeErr := jwt.Decode([]byte(tokenstring), &claimsDecoded, []byte("azz"))
	if decodeErr != nil {
		glog.Glog(LogF, "401: Not Authorized")
		return errors.New("401: Not Authorized")
	}

	exp := claimsDecoded["exp"].(float64)
	exp1, _ := strconv.ParseFloat(fmt.Sprintf("%v", time.Now().Unix()+0), 64)

	if (exp - exp1) < 0 {
		glog.Glog(LogF, fmt.Sprintf("401: Not Authorized AccessToken Expired %v %v %v ,Please login", exp, exp1, (exp-exp1)))
		return errors.New(fmt.Sprintf("401: Not Authorized AccessToken Expired %v %v %v ,Please login", exp, exp1, (exp - exp1)))
	}
	return nil
}

type MyWork struct {
	Id   string "id"
	Name string "The Name of a process"
	WP   *workpool.WorkPool
	Mmw  *MetaMyWork
}

func (workPool *MyWork) DoWork(workRoutine int) {
	glog.Glog(LogF, fmt.Sprintf("%s", workPool.Name))
	glog.Glog(LogF, fmt.Sprintf("*******> WR: %d  QW: %d  AR: %d", workRoutine, workPool.WP.QueuedWork(), workPool.WP.ActiveRoutines()))

	nmq := NewMstMgr(workPool.Mmw.FlowId, workPool.Mmw.ApiServerIp, workPool.Mmw.ApiServerPort, fmt.Sprintf("m%v", workRoutine), workPool.Mmw.HomeDir)
	nmq.MainMgr()
}
