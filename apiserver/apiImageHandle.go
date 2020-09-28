package apiserver

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/db"
	"github.com/lhzd863/autoflow/glog"
	"github.com/lhzd863/autoflow/module"
	"github.com/lhzd863/autoflow/util"

	uuid "github.com/satori/go.uuid"
)

type ResponseResourceImage struct {
	sync.Mutex
	Conf *module.MetaApiServerBean
}

func NewResponseResourceImage(conf *module.MetaApiServerBean) *ResponseResourceImage {
	return &ResponseResourceImage{Conf: conf}
}

func (rrs *ResponseResourceImage) ImageAddHandler(request *restful.Request, response *restful.Response) {
	reqParams, err := url.ParseQuery(request.Request.URL.RawQuery)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, "parse url parameter err.", nil)
		return
	}

	username, err := util.JwtAccessTokenUserName(fmt.Sprint(reqParams["accesstoken"][0]), conf.JwtKey)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("accesstoken parse username err.%v", err), nil)
		return
	}

	p := new(module.MetaParaImageAddBean)
	err = request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}

	if !util.FileExist(p.DbStore) {
		glog.Glog(LogF, p.DbStore+" not exists.")
		util.ApiResponse(response.ResponseWriter, 700, p.DbStore+" not exists.", nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()

	m := new(module.MetaJobImageBean)

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.User = username

	u1 := uuid.Must(uuid.NewV4(), nil)
	imageid := fmt.Sprint(u1)
	m.ImageId = imageid
	m.Tag = p.Tag
	m.DbStore = p.DbStore
	m.Description = p.Description
	m.Enable = p.Enable

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(imageid, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}

	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceImage) ImageUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaImageUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.ImageId) < 1 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()
	fb0 := bt.Get(p.ImageId)
	fb := new(module.MetaJobImageBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Tag = p.Tag
	fb.DbStore = p.DbStore
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.ImageId, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceImage) ImageRemoveHandler(request *restful.Request, response *restful.Response) {

	imagebean := new(module.MetaParaImageRemoveBean)
	err := request.ReadEntity(&imagebean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()

	err = bt.Remove(imagebean.ImageId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceImage) ImageListHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			masb := new(module.MetaJobImageBean)
			err := json.Unmarshal([]byte(v1.(string)), &masb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, masb)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceImage) ImageGetHandler(request *restful.Request, response *restful.Response) {
	imagebean := new(module.MetaParaImageGetBean)
	err := request.ReadEntity(&imagebean)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_IMAGE)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(imagebean.ImageId)
	if ib != nil {
		masb := new(module.MetaJobImageBean)
		err := json.Unmarshal([]byte(ib.(string)), &masb)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, masb)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}
