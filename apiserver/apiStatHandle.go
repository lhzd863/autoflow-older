package apiserver

import (
        "encoding/json"
        "fmt"
        "sync"

        "github.com/emicklei/go-restful"

        "github.com/lhzd863/autoflow/internal/db"
        "github.com/lhzd863/autoflow/internal/glog"
        "github.com/lhzd863/autoflow/internal/module"
        "github.com/lhzd863/autoflow/internal/util"
)

type ResponseResourceStat struct {
        sync.Mutex
}

func NewResponseResourceStat() *ResponseResourceStat {
        return &ResponseResourceStat{}
}

func (rrs *ResponseResourceStat) SystemStatJobStatusCntHandler(request *restful.Request, response *restful.Response) {
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
        retmap := make(map[string]int)
        retmap["Ready"] = 0
        retmap["Pending"] = 0
        retmap["Submit"] = 0
        retmap["Go"] = 0
        retmap["Succ"] = 0
        retmap["Fail"] = 0
        retmap["Running"] = 0
        retlst := make([]interface{}, 0)
        for _, v := range strlist {
                for _, v1 := range v.(map[string]interface{}) {
                        m := new(module.MetaJobBean)
                        err := json.Unmarshal([]byte(v1.(string)), &m)
                        if err != nil {
                                glog.Glog(LogF, fmt.Sprint(err))
                                continue
                        }
                        if _, ok := retmap[m.Status]; ok {
                           t := retmap[m.Status]
                           retmap[m.Status] = t+1
                        }else {
                           retmap[m.Status] = 1
                        }
                }
        }

        n := new(module.MetaSystemStatJobStatusBean)
        n.FlowId = p.FlowId
        n.Ready = fmt.Sprint(retmap["Ready"])
        n.Pending = fmt.Sprint(retmap["Pending"])
        n.Submit = fmt.Sprint(retmap["Submit"])
        n.Go = fmt.Sprint(retmap["Go"])
        n.Succ = fmt.Sprint(retmap["Succ"])
        n.Fail = fmt.Sprint(retmap["Fail"])
        n.Running = fmt.Sprint(retmap["Running"])
        retlst = append(retlst, n)
        util.ApiResponse(response.ResponseWriter, 200, "", retlst)       
}

