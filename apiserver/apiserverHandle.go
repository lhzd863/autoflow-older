package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/gproto"
	"github.com/lhzd863/autoflow/internal/jwt"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"

	"github.com/satori/go.uuid"
)

// JobResource is the REST layer to the
type ResponseResource struct {
	// normally one would use DAO (data access object)
	//Data map[string]interface{}
	sync.Mutex
}

// WebService creates a new service that can handle REST requests for User resources.
func (rrs ResponseResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1").
		Consumes("*/*").
		Produces(restful.MIME_JSON, restful.MIME_JSON) // you can specify this per route as well

	tags := []string{"system"}

	ws.Route(ws.GET("/health").To(rrs.HealthHandler).
		// docs
		Doc("Health").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowBean{}).
		Writes(ResponseResource{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/login1").To(rrs.LoginHandler).
		// docs
		Doc("login info").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(ResponseResource{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/login").To(rru.SystemUserTokenHandler).
		// docs
		Doc("login info").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(module.MetaParaSystemUserTokenBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-user"}
	ws.Route(ws.POST("/system/user/add").To(rru.SystemUserAddHandler).
		// docs
		Doc("新增用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/get").To(rru.SystemUserGetHandler).
		// docs
		Doc("获取用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/rm").To(rru.SystemUserRemoveHandler).
		// docs
		Doc("删除用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/update").To(rru.SystemUserUpdateHandler).
		// docs
		Doc("更新用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/ls").To(rru.SystemUserListHandler).
		// docs
		Doc("用户列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-role"}
	ws.Route(ws.POST("/system/role/add").To(rru.SystemRoleAddHandler).
		// docs
		Doc("新增角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/get").To(rru.SystemRoleGetHandler).
		// docs
		Doc("获取角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/rm").To(rru.SystemRoleRemoveHandler).
		// docs
		Doc("删除角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/update").To(rru.SystemRoleUpdateHandler).
		// docs
		Doc("更新角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/ls").To(rru.SystemRoleListHandler).
		// docs
		Doc("列表角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-right"}
	ws.Route(ws.POST("/system/right/add").To(rru.SystemRightAddHandler).
		// docs
		Doc("新增权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRightAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/right/get").To(rru.SystemRightGetHandler).
		// docs
		Doc("获取权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRightGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/right/rm").To(rru.SystemRightRemoveHandler).
		// docs
		Doc("删除权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRightRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/right/update").To(rru.SystemRightUpdateHandler).
		// docs
		Doc("更新权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRightUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/right/ls").To(rru.SystemRightListHandler).
		// docs
		Doc("列表权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-role-right"}
	ws.Route(ws.POST("/system/role/right/add").To(rru.SystemRoleRightAddHandler).
		// docs
		Doc("新增角色权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleRightAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/right/get").To(rru.SystemRoleRightGetHandler).
		// docs
		Doc("获取角色权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleRightGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/right/rm").To(rru.SystemRoleRightRemoveHandler).
		// docs
		Doc("删除角色权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleRightRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/right/update").To(rru.SystemRoleRightUpdateHandler).
		// docs
		Doc("更新角色权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRoleRightUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/role/right/ls").To(rru.SystemRoleRightListHandler).
		// docs
		Doc("列表角色权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-user-role"}
	ws.Route(ws.POST("/system/user/role/add").To(rru.SystemUserRoleAddHandler).
		// docs
		Doc("新增用户角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserRoleAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/role/get").To(rru.SystemUserRoleGetHandler).
		// docs
		Doc("获取用户角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserRoleGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/role/rm").To(rru.SystemUserRoleRemoveHandler).
		// docs
		Doc("删除用户角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserRoleRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/role/update").To(rru.SystemUserRoleUpdateHandler).
		// docs
		Doc("更新用户角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemUserRoleUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/user/role/ls").To(rru.SystemUserRoleListHandler).
		// docs
		Doc("列表用户角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", ResponseResource{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-parameter"}
	ws.Route(ws.POST("/system/parameter/add").To(rrs.SystemParameterAddHandler).
		// docs
		Doc("新增系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/rm").To(rrs.SystemParameterRemoveHandler).
		// docs
		Doc("删除系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/get").To(rrs.SystemParameterGetHandler).
		// docs
		Doc("获取系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/ls").To(rrs.SystemParameterListHandler).
		// docs
		Doc("列表系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/update").To(rrs.SystemParameterUpdateHandler).
		// docs
		Doc("更新系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system"}
	ws.Route(ws.POST("/sys/lsport").To(rrs.SysListPortHandler).
		// docs
		Doc("list port").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"image"}
	ws.Route(ws.POST("/image/add").To(rrs.ImageAddHandler).
		// docs
		Doc("新增镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/rm").To(rrs.ImageRemoveHandler).
		// docs
		Doc("删除镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/ls").To(rrs.ImageListHandler).
		// docs
		Doc("列表镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/get").To(rrs.ImageGetHandler).
		// docs
		Doc("获取镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/update").To(rrs.ImageUpdateHandler).
		// docs
		Doc("更新镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"instance"}
	ws.Route(ws.POST("/instance/create").To(rrs.InstanceCreateHandler).
		// docs
		Doc("实例创建").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaJobFlowBean{ImageId: "", ProcessNum: ""}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/start").To(rrs.InstanceStartHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/stop").To(rrs.InstanceStopHandler).
		// docs
		Doc("实例停止").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/ls").To(rrs.InstanceListHandler).
		// docs
		Doc("实例列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/ls/status").To(rrs.InstanceListStatusHandler).
		// docs
		Doc("实例状态列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/rm").To(rrs.InstanceRemoveHandler).
		// docs
		Doc("实例列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow"}
	ws.Route(ws.POST("/mst/instance/start").To(rrs.MstInstanceStartHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowStartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/instance/stop").To(rrs.MstInstanceStopHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowStopBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-routine"}
	ws.Route(ws.POST("/flow/routine/add").To(rrs.FlowRoutineAddHandler).
		// docs
		Doc("实例新增处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowRoutineAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/routine/sub").To(rrs.FlowRoutineSubHandler).
		// docs
		Doc("实例线程删除").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowRoutineSubBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/routine/status").To(rrs.FlowRoutineStatusHandler).
		// docs
		Doc("流获取").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine"}
	ws.Route(ws.POST("/mst/flow/routine/start").To(rrs.MstFlowRoutineStartHandler).
		// docs
		Doc("指定mst实例新增处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/stop").To(rrs.MstFlowRoutineStopHandler).
		// docs
		Doc("指定mst实例停止处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow"}
	ws.Route(ws.POST("/flow/ls").To(rrs.FlowListHandler).
		// docs
		Doc("列表实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/get").To(rrs.FlowGetHandler).
		// docs
		Doc("获取实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/update").To(rrs.FlowUpdateHandler).
		// docs
		Doc("更新实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/status/update").To(rrs.FlowStatusUpdateHandler).
		// docs
		Doc("更新实例状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowStatusUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst"}
	ws.Route(ws.POST("/mst/flow/rm").To(rrs.MstFlowRemoveHandler).
		// docs
		Doc("删除实例MST映射").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-parameter"}
	ws.Route(ws.POST("/flow/parameter/add").To(rrs.FlowParameterAddHandler).
		// docs
		Doc("新增实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/ls").To(rrs.FlowParameterListHandler).
		// docs
		Doc("列表实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/get").To(rrs.FlowParameterGetHandler).
		// docs
		Doc("获取实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/rm").To(rrs.FlowParameterRemoveHandler).
		// docs
		Doc("删除实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/update").To(rrs.FlowParameterUpdateHandler).
		// docs
		Doc("更新实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-pool"}
	ws.Route(ws.POST("/job/pool/add").To(rrs.JobPoolAddHandler).
		// docs
		Doc("新增作业到系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/rm").To(rrs.JobPoolRemoveHandler).
		// docs
		Doc("删除作业从系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/get").To(rrs.JobPoolGetHandler).
		// docs
		Doc("获取作业从系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/ls").To(rrs.JobPoolListHandler).
		// docs
		Doc("列表在系统作业池中作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker"}
	ws.Route(ws.POST("/worker/heart/add").To(rrs.WorkerHeartAddHandler).
		// docs
		Doc("注册执行节点").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/rm").To(rrs.WorkerHeartRemoveHandler).
		// docs
		Doc("节点删除").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/ls").To(rrs.WorkerHeartListHandler).
		// docs
		Doc("节点列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/get").To(rrs.WorkerHeartGetHandler).
		// docs
		Doc("节点获取").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker-cnt"}
	ws.Route(ws.POST("/worker/cnt/add").To(rrs.WorkerCntAddHandler).
		// docs
		Doc("节点作业数增加").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaWorkerHeartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/cnt/exec").To(rrs.WorkerExecCntHandler).
		// docs
		Doc("节点执行作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job"}
	ws.Route(ws.POST("/flow/job/ls").To(rrs.FlowJobListHandle).
		// docs
		Doc("列表实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/get").To(rrs.FlowJobGetHandle).
		// docs
		Doc("获取实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/add").To(rrs.FlowJobAddHandle).
		// docs
		Doc("新增实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/update").To(rrs.FlowJobUpdateHandle).
		// docs
		Doc("更新实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/rm").To(rrs.FlowJobRemoveHandle).
		// docs
		Doc("删除实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-status-get"}
	ws.Route(ws.POST("/flow/job/status/get/pending").To(rrs.FlowJobStatusGetPendingHandle).
		// docs
		Doc("获取Pending状态实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusGetPendingBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/get/go").To(rrs.FlowJobStatusGetGoHandle).
		// docs
		Doc("获取Go状态实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusGetGoBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-heart"}
	ws.Route(ws.POST("/mst/heart/add").To(rrs.MstHeartAddHandler).
		// docs
		Doc("新增Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/ls").To(rrs.MstHeartListHandler).
		// docs
		Doc("列表Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/rm").To(rrs.MstHeartRemoveHandler).
		// docs
		Doc("删除Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/get").To(rrs.MstHeartGetHandler).
		// docs
		Doc("获取Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine-heart"}
	ws.Route(ws.POST("/mst/flow/routine/heart/add").To(rrs.MstFlowRoutineHeartAddHandler).
		// docs
		Doc("新增Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/ls").To(rrs.MstFlowRoutineHeartListHandler).
		// docs
		Doc("列表Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/rm").To(rrs.MstFlowRoutineHeartRemoveHandler).
		// docs
		Doc("删除Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/get").To(rrs.MstFlowRoutineHeartGetHandler).
		// docs
		Doc("获取Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-dependency"}
	ws.Route(ws.POST("/flow/job/dependency").To(rrs.FlowJobDependencyHandler).
		// docs
		Doc("流作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/ls").To(rrs.FlowJobDependencyListHandler).
		// docs
		Doc("列表作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/rm").To(rrs.FlowJobDependencyRemoveHandler).
		// docs
		Doc("删除作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/get").To(rrs.FlowJobDependencyGetHandler).
		// docs
		Doc("获取作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/update").To(rrs.FlowJobDependencyUpdateHandler).
		// docs
		Doc("更新作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/add").To(rrs.FlowJobDependencyAddHandler).
		// docs
		Doc("新增作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-timewindow"}
	ws.Route(ws.POST("/flow/job/timewindow").To(rrs.FlowJobTimeWindowHandler).
		// docs
		Doc("实例作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/add").To(rrs.FlowJobTimeWindowAddHandler).
		// docs
		Doc("新增作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/rm").To(rrs.FlowJobTimeWindowRemoveHandler).
		// docs
		Doc("删除作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/ls").To(rrs.FlowJobTimeWindowListHandler).
		// docs
		Doc("列表作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/get").To(rrs.FlowJobTimeWindowGetHandler).
		// docs
		Doc("获取作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/update").To(rrs.FlowJobTimeWindowUpdateHandler).
		// docs
		Doc("更新作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-stream"}
	ws.Route(ws.POST("/flow/job/stream/job").To(rrs.FlowJobStreamJobHandler).
		// docs
		Doc("触发作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamJobBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/job/get").To(rrs.FlowJobStreamJobGetHandler).
		// docs
		Doc("触发作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamJobListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/add").To(rrs.FlowJobStreamAddHandler).
		// docs
		Doc("新增作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/rm").To(rrs.FlowJobStreamRemoveHandler).
		// docs
		Doc("删除作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/ls").To(rrs.FlowJobStreamListHandler).
		// docs
		Doc("列表作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/get").To(rrs.FlowJobStreamGetHandler).
		// docs
		Doc("获取作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/update").To(rrs.FlowJobStreamUpdateHandler).
		// docs
		Doc("更新作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-status-update"}
	ws.Route(ws.POST("/flow/job/status/update/submit").To(rrs.FlowJobStatusUpdateSubmitHandler).
		// docs
		Doc("更新作业状态为Submit").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateSubmitBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/go").To(rrs.FlowJobStatusUpdateGoHandler).
		// docs
		Doc("更新作业状态为Go").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateGoBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/pending").To(rrs.FlowJobStatusUpdatePendingHandler).
		// docs
		Doc("更新作业状态为Pending").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdatePendingBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/start").To(rrs.FlowJobStatusUpdateStartHandler).
		// docs
		Doc("更新作业开始运行相关信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateStartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/end").To(rrs.FlowJobStatusUpdateEndHandler).
		// docs
		Doc("更新作业结束运行相关信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateEndBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-cmd"}
	ws.Route(ws.POST("/flow/job/cmd/add").To(rrs.FlowJobCmdAddHandler).
		// docs
		Doc("新增作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/getall").To(rrs.FlowJobCmdGetAllHandler).
		// docs
		Doc("获取作业所有脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Reads(module.MetaParaFlowJobCmdGetAllBean{}).
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/get").To(rrs.FlowJobCmdGetHandler).
		// docs
		Doc("获取作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Reads(module.MetaParaFlowJobCmdGetBean{}).
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/ls").To(rrs.FlowJobCmdListHandler).
		// docs
		Doc("列表作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/update").To(rrs.FlowJobCmdUpdateHandler).
		// docs
		Doc("更新作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/rm").To(rrs.FlowJobCmdRemoveHandler).
		// docs
		Doc("删除作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-parameter"}
	ws.Route(ws.POST("/flow/job/parameter/add").To(rrs.FlowJobParameterAddHandler).
		// docs
		Doc("新增作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/get").To(rrs.FlowJobParameterGetHandler).
		// docs
		Doc("获取作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/getall").To(rrs.FlowJobParameterGetAllHandler).
		// docs
		Doc("列表系统实例作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterGetAllBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/ls").To(rrs.FlowJobParameterListHandler).
		// docs
		Doc("列表作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/update").To(rrs.FlowJobParameterUpdateHandler).
		// docs
		Doc("更新作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/rm").To(rrs.FlowJobParameterRemoveHandler).
		// docs
		Doc("删除作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-log"}
	ws.Route(ws.POST("/flow/job/log/add").To(rrs.FlowJobLogAddHandler).
		// docs
		Doc("新增作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/get").To(rrs.FlowJobLogGetHandler).
		// docs
		Doc("获取作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/ls").To(rrs.FlowJobLogListHandler).
		// docs
		Doc("列表作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/rm").To(rrs.FlowJobLogRemoveHandler).
		// docs
		Doc("删除作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/append").To(rrs.FlowJobLogAppendHandler).
		// docs
		Doc("追加作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine-job-running-heart"}
	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/add").To(rrs.MstFlowRoutineJobRunningHeartAddHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/rm").To(rrs.MstFlowRoutineJobRunningHeartRemoveHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/get").To(rrs.MstFlowRoutineJobRunningHeartGetHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/ls").To(rrs.MstFlowRoutineJobRunningHeartListHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker-routine-job-running-heart"}
	ws.Route(ws.POST("/worker/routine/job/running/heart/ls").To(rrs.WorkerRoutineJobRunningHeartListHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/get").To(rrs.WorkerRoutineJobRunningHeartGetHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/rm").To(rrs.WorkerRoutineJobRunningHeartRemoveHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/add").To(rrs.WorkerRoutineJobRunningHeartAddHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-ring"}
	ws.Route(ws.POST("/system/ring/pending/ls").To(rrs.SystemRingPendingListHandler).
		// docs
		Doc("列表Pending Ring").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/go/ls").To(rrs.SystemRingGoListHandler).
		// docs
		Doc("列表Go Ring").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/pending/rm").To(rrs.SystemRingPendingRemoveHandler).
		// docs
		Doc("删除Pending Ring信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRingPendingRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/go/rm").To(rrs.SystemRingGoRemoveHandler).
		// docs
		Doc("删除Go Ring信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRingGoRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	return ws
}

func (rrs *ResponseResource) HealthHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) LoginHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
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

func (rrs *ResponseResource) SystemParameterListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SystemParameterGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SystemParameterUpdateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SystemParameterRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SystemParameterAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) ImageAddHandler(request *restful.Request, response *restful.Response) {
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

	u1 := uuid.Must(uuid.NewV4())
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

func (rrs *ResponseResource) ImageUpdateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) ImageRemoveHandler(request *restful.Request, response *restful.Response) {

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

func (rrs *ResponseResource) ImageListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowUpdateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) ImageGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) InstanceCreateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) InstanceStartHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstInstanceStartHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) InstanceStopHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstInstanceStopHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineStartHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineStopHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) InstanceListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) InstanceListStatusHandler(request *restful.Request, response *restful.Response) {

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

func (rrs *ResponseResource) flowInfo(flowid string) (*module.MetaJobFlowBean, error) {
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

func (rrs *ResponseResource) InstanceRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowRoutineStatusHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowStatusUpdateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowRoutineAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowRoutineSubHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowParameterListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowParameterGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowParameterUpdateHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowParameterRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowParameterAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SysListPortHandler(request *restful.Request, response *restful.Response) {
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
			ism, _ := rrs.IsExpiredMst(mmf.MstId)
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

func (rrs *ResponseResource) JobPoolAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) JobPoolGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) JobPoolListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) JobPoolRemoveHandler(request *restful.Request, response *restful.Response) {

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

func (rrs *ResponseResource) getMstPort(pnum int) []*module.MetaMstHeartBean {
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

func (rrs *ResponseResource) getMstInfo(mstid string, pnum int) (*module.MetaMstHeartBean, error) {
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

func (rrs *ResponseResource) getMstHeart() []*module.MetaMstHeartBean {
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

func getMstHeart1(mstid string) ([]module.MetaMstHeartBean, error) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_MASTER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	flag := 0
	retlst := make([]module.MetaMstHeartBean, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mmhb := new(module.MetaMstHeartBean)
			err := json.Unmarshal([]byte(v1.(string)), &mmhb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if mstid != mmhb.MstId {
				continue
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			ise, _ := util.IsExpired(mmhb.UpdateTime, timeStr, 300)
			if ise {
				glog.Glog(LogF, fmt.Sprintf("%v timeout %v,%v.", mmhb.MstId, mmhb.Ip, mmhb.Port))
				continue
			}
			retlst = append(retlst, *mmhb)
			flag = 1
		}
	}

	if flag == 0 {
		return retlst, errors.New(fmt.Sprintf("%v mst ip mapping no port,wait for next time.", mstid))
	}
	return retlst, nil
}

//func saveMstFlow() (*module.MetaInstancePortBean, error) {
func (rrs *ResponseResource) saveMstFlow(mmfb *module.MetaMstFlowBean) error {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_FLOW_MASTER)
	defer bt.Close()

	jsonstr, _ := json.Marshal(mmfb)
	err := bt.Set(mmfb.MstId+"."+mmfb.FlowId, string(jsonstr))
	if err != nil {
		return errors.New(fmt.Sprintf("data in db update error.%v", err))
	}
	return nil
}

func (rrs *ResponseResource) getFlowPort(flowid string) ([]module.MetaMstFlowBean, error) {
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

func (rrs *ResponseResource) getMstFlowPort(flowid string, mstid string) []*module.MetaMstFlowRoutineHeartBean {
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

func (rrs *ResponseResource) getMstFlowRoutinePort(flowid string, mstid string, routineid string) []*module.MetaMstFlowRoutineHeartBean {
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
			if mmf.FlowId != flowid || mmf.MstId != mstid || mmf.RoutineId != routineid {
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

func (rrs *ResponseResource) IsExpiredMst(mstid string) (bool, error) {
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

func (rrs *ResponseResource) WorkerHeartAddHandler(request *restful.Request, response *restful.Response) {
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
	m := new(module.MetaWorkerHeartBean)
	m.Id = p.Id
	m.WorkerId = p.WorkerId
	m.Ip = p.Ip
	m.Port = p.Port
	m.MaxCnt = p.MaxCnt
	m.RunningCnt = p.RunningCnt
	m.CurrentCnt = p.CurrentCnt
	m.StartTime = p.StartTime
	m.Duration = p.Duration
	rrs.Lock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	rrs.Unlock()
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) WorkerHeartRemoveHandler(request *restful.Request, response *restful.Response) {
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
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	rrs.Unlock()
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) WorkerHeartListHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaWorkerHeartBean)
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

func (rrs *ResponseResource) WorkerHeartGetHandler(request *restful.Request, response *restful.Response) {
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
	rrs.Unlock()
	if m != nil {
		v := new(module.MetaWorkerHeartBean)
		err := json.Unmarshal([]byte(m.(string)), &v)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) WorkerCntAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaWorkerHeartBean)
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
	fb := new(module.MetaWorkerHeartBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.CurrentCnt = p.CurrentCnt
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

func (rrs *ResponseResource) WorkerExecCntHandler(request *restful.Request, response *restful.Response) {
	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_WORKER_HEART)
	defer bt.Close()

	strlist := bt.Scan()
	//bt.Close()

	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for k1, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaWorkerHeartBean)
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
			if est-sst > 600 {
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
			currentcnt, err := strconv.Atoi(m.CurrentCnt)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("conv currentcnt fail.%v", err))
				continue
			}
			if maxcnt <= runningcnt+currentcnt {
				glog.Glog(LogF, "max cnt gt running job cnt.")
				continue
			}
			retlst = append(retlst, m)
		}
	}
	bt.Close()
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobListHandle(request *restful.Request, response *restful.Response) {
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
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobBean)
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

func (rrs *ResponseResource) FlowJobGetHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	m := bt.Get(p.Sys + "." + p.Job)
	if m != nil {
		r := new(module.MetaJobBean)
		err := json.Unmarshal([]byte(m.(string)), &r)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, r)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobAddHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 != nil {
		glog.Glog(LogF, fmt.Sprintf("%v %v has exists.", p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v %v has exists.", p.Sys, p.Job), nil)
		return
	}
	mj := new(module.MetaJobBean)
	mj.Sys = p.Sys
	mj.Job = p.Job
	mj.Enable = p.Enable
	mj.TimeWindow = p.TimeWindow
	mj.RetryCnt = p.RetryCnt
	mj.Alert = p.Alert
	mj.TimeTrigger = p.TimeTrigger
	mj.JobType = p.JobType
	mj.Frequency = p.Frequency
	mj.CheckBatStatus = p.CheckBatStatus
	mj.Priority = p.Priority
	jsonstr, _ := json.Marshal(mj)
	err = bt.Set(mj.Sys+"."+mj.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobUpdateHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 != nil {
		mj := new(module.MetaJobBean)
		err := json.Unmarshal([]byte(fb0.(string)), &mj)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		mj.Enable = p.Enable
		mj.TimeWindow = p.TimeWindow
		mj.RetryCnt = p.RetryCnt
		mj.Alert = p.Alert
		mj.TimeTrigger = p.TimeTrigger
		mj.JobType = p.JobType
		mj.Frequency = p.Frequency
		mj.CheckBatStatus = p.CheckBatStatus
		mj.Priority = p.Priority
		mj.Status = p.Status
		jsonstr, _ := json.Marshal(mj)
		err = bt.Set(mj.Sys+"."+mj.Job, string(jsonstr))
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
			return
		}
	} else {
		glog.Glog(LogF, fmt.Sprintf("flow %v job %v %v no store data.", p.FlowId, p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("flow %v job %v %v no store data.", p.FlowId, p.Sys, p.Job), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobRemoveHandle(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	bt.Remove(p.Sys + "," + p.Job)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusGetPendingHandle(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusGetPendingBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	c, bid := rrs.CurrentStatusPendingOffset(len(strlist))
	retlst := make([]interface{}, 0)
	if bid == -1 {
		util.ApiResponse(response.ResponseWriter, 200, fmt.Sprintf("all ringid is working."), retlst)
		return
	} else {
		r := new(module.MetaRingPendingOffsetBean)
		r.Id = jobpara.Id
		r.RingId = fmt.Sprint(bid)
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		r.CreateTime = timeStr
		ringPendingSpool.Add(jobpara.Id, r)
	}
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobbn := new(module.MetaJobBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobbn)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobbn.Status != jobpara.Status && jobpara.Status != "ALL" {
				continue
			}
			if (c && jobpara.IsHash == "0") || jobpara.IsHash == "1" {
				jobnodeid := statusPendingHashRing.Get(jobbn.Job).Id
				glog.Glog(LogF, fmt.Sprintf("job hash info:%v %v", jobbn.Job, jobnodeid))
				if jobnodeid != bid {
					glog.Glog(LogF, fmt.Sprintf("local node id %v ,job %v  mapping node id %v is not local node.", bid, jobbn.Job, jobnodeid))
					continue
				}
			}
			retlst = append(retlst, jobbn)
		}
	}
	if len(retlst) == 0 {
		ringPendingSpool.Remove(jobpara.Id)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusGetGoHandle(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusGetGoBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	c, bid := rrs.CurrentStatusGoOffset(len(strlist))
	retlst := make([]interface{}, 0)
	if bid == -1 {
		util.ApiResponse(response.ResponseWriter, 200, "all ringid is working.", retlst)
		return
	} else {
		r := new(module.MetaRingPendingOffsetBean)
		r.Id = jobpara.Id
		r.RingId = fmt.Sprint(bid)
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		r.CreateTime = timeStr
		ringGoSpool.Add(jobpara.Id, r)
	}
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobbn := new(module.MetaJobBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobbn)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobbn.Status != jobpara.Status && jobpara.Status != "ALL" {
				continue
			}
			if (c && jobpara.IsHash == "0") || jobpara.IsHash == "1" {
				jobnodeid := statusGoHashRing.Get(jobbn.Job).Id
				glog.Glog(LogF, fmt.Sprintf("job hash info:%v %v", jobbn.Job, jobnodeid))
				if jobnodeid != bid {
					glog.Glog(LogF, fmt.Sprintf("local node id %v ,job %v  mapping node id %v is not local node.", bid, jobbn.Job, jobnodeid))
					continue
				}
			}
			retlst = append(retlst, jobbn)
		}
	}
	if len(retlst) == 0 {
		ringGoSpool.Remove(jobpara.Id)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) MstHeartAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstHeartListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstHeartGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstHeartRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineHeartAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineHeartListHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineHeartGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineHeartRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) CurrentStatusPendingOffset(size int) (bool, int) {
	rrs.Lock()
	defer rrs.Unlock()
	t := -1
	c := false
	v := statusPendingOffset
	for i := 0; i < len(statusPendingHashRing.Resources); i++ {
		if statusPendingOffset == len(statusPendingHashRing.Resources)-1 {
			statusPendingOffset = 0
		} else {
			statusPendingOffset += 1
		}
		if rrs.IsExistPendingRingId(fmt.Sprint(v)) {
			v = statusPendingOffset
		} else {
			t = v
			break
		}
	}
	if size > 300 {
		c = true
	}
	return c, t
}

func (rrs *ResponseResource) IsExistPendingRingId(ringid string) bool {
	for k := range ringPendingSpool.MemMap {
		v := ringPendingSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		if v.RingId == ringid {
			loc, _ := time.LoadLocation("Local")
			timeLayout := "2006-01-02 15:04:05"
			stheTime, _ := time.ParseInLocation(timeLayout, v.CreateTime, loc)
			sst := stheTime.Unix()
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
			est := etheTime.Unix()
			if est-sst > 300 {
				ringPendingSpool.Remove(k)
				break
			}
			return true
		}
	}
	return false
}

func (rrs *ResponseResource) CurrentStatusGoOffset(size int) (bool, int) {
	rrs.Lock()
	defer rrs.Unlock()
	t := -1
	c := false
	v := statusGoOffset
	for i := 0; i < len(statusGoHashRing.Resources); i++ {
		if statusGoOffset == len(statusGoHashRing.Resources)-1 {
			statusGoOffset = 0
		} else {
			statusGoOffset += 1
		}
		if rrs.IsExistGoRingId(fmt.Sprint(v)) {
			v = statusGoOffset
		} else {
			t = v
			break
		}
	}

	if size > 300 {
		c = true
	}
	return c, t
}

func (rrs *ResponseResource) IsExistGoRingId(ringid string) bool {
	for k := range ringGoSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		if v.RingId == ringid {
			loc, _ := time.LoadLocation("Local")
			timeLayout := "2006-01-02 15:04:05"
			stheTime, _ := time.ParseInLocation(timeLayout, v.CreateTime, loc)
			sst := stheTime.Unix()
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			etheTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
			est := etheTime.Unix()
			if est-sst > 60 {
				ringGoSpool.Remove(k)
				break
			}
			return true
		}
	}
	return false
}

func (rrs *ResponseResource) SystemRingGoListHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	for k := range ringGoSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingGoOffsetBean)
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) SystemRingGoRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) SystemRingPendingListHandler(request *restful.Request, response *restful.Response) {
	retlst := make([]interface{}, 0)
	for k := range ringPendingSpool.MemMap {
		v := ringGoSpool.MemMap[k].(module.MetaRingPendingOffsetBean)
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) SystemRingPendingRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) instanceHomeDir(flowid string) (string, error) {

	return conf.HomeDir, nil
}

func (rrs *ResponseResource) flowDbFile(flowid string) (string, error) {
	f := conf.HomeDir + "/" + flowid + "/" + flowid + ".db"
	return f, nil
}

func (rrs *ResponseResource) FlowJobDependencyHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobDependencyBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			jobdb := new(module.MetaJobDependencyBean)
			err := json.Unmarshal([]byte(v1.(string)), &jobdb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if jobdb.Enable != "1" {
				continue
			}
			if jobdb.Sys == jobpara.Sys && jobdb.Job == jobpara.Job {
				bt1 := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
				defer bt1.Close()
				jobn := bt1.Get(jobdb.DependencySys + "." + jobdb.DependencyJob)
				bt1.Close()
				if jobn != nil {
					jobbn := new(module.MetaJobBean)
					err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
					if err != nil {
						glog.Glog(LogF, fmt.Sprint(err))
					}
					if jobbn.Status == util.STATUS_AUTO_SUCC || jobbn.Enable != "1" || jobdb.Enable != "1" {
						continue
					}
					retlst = append(retlst, jobdb)
				}
			}
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobDependencyListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobDependencyBean)
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

func (rrs *ResponseResource) FlowJobDependencyGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	if ib != nil {
		m := new(module.MetaJobDependencyBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobDependencyUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()
	fb0 := bt.Get(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	fb := new(module.MetaJobDependencyBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.Sys+"."+p.Job+"."+p.DependencySys+"."+p.DependencyJob, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobDependencyRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	err = bt.Remove(p.Sys + "." + p.Job + "." + p.DependencySys + "." + p.DependencyJob)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobDependencyAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobDependencyAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.DependencySys) == 0 || len(p.DependencyJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_DEPENDENCY)
	defer bt.Close()

	jsonstr, _ := json.Marshal(p)
	err = bt.Set(p.Sys+"."+p.Job+"."+p.DependencySys+"."+p.DependencyJob, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamJobHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamJobBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	fb0 := bt.Get(p.Sys + "." + p.Job)
	if fb0 == nil {
		glog.Glog(LogF, fmt.Sprintf("%v.%v not exists.", p.Sys, p.Job))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v.%v not exists.", p.Sys, p.Job), nil)
		return
	}
	fb := new(module.MetaJobBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}

	if fb.Status != util.STATUS_AUTO_SUCC && fb.Status != util.STATUS_AUTO_READY {
		glog.Glog(LogF, fmt.Sprintf("%v.%v status %v not equal %v or %v .", p.Sys, p.Job, fb.Status, util.STATUS_AUTO_SUCC, util.STATUS_AUTO_READY))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v.%v status %v not equal %v or %v .", p.Sys, p.Job, fb.Status, util.STATUS_AUTO_SUCC, util.STATUS_AUTO_READY), nil)
		return
	}

	fb.Status = util.STATUS_AUTO_PENDING
	fb.EndTime = ""
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fb.StartTime = timeStr
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Sys+"."+fb.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamJobGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamJobListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobStreamBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if m.StreamSys == p.Sys && m.StreamJob == p.Job {
				retlst = append(retlst, m)
			}
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobStreamBean)
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

func (rrs *ResponseResource) FlowJobStreamGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	if ib != nil {
		m := new(module.MetaJobStreamBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()
	fb0 := bt.Get(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	fb := new(module.MetaJobStreamBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Description = p.Description
	fb.Enable = p.Enable

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.StreamSys+"."+p.StreamJob+"."+p.Sys+"."+p.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	err = bt.Remove(p.StreamSys + "." + p.StreamJob + "." + p.Sys + "." + p.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStreamAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobStreamAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.StreamSys) == 0 || len(p.StreamJob) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_STREAM)
	defer bt.Close()

	m := new(module.MetaJobStreamBean)
	m.StreamSys = p.StreamSys
	m.StreamJob = p.StreamJob
	m.Sys = p.Sys
	m.Job = p.Job
	m.Description = p.Description
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.StreamSys+"."+m.StreamJob+"."+m.Sys+"."+m.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobTimeWindowListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobTimeWindowBean)
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

func (rrs *ResponseResource) FlowJobTimeWindowGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job)
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

func (rrs *ResponseResource) FlowJobTimeWindowUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowUpdateBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()
	fb0 := bt.Get(p.Sys + "." + p.Job)
	fb := new(module.MetaJobTimeWindowBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Allow = p.Allow
	fb.StartHour = p.StartHour
	fb.EndHour = p.EndHour
	fb.Enable = p.Enable
	fb.Description = p.Description

	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(p.Sys+"."+p.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobTimeWindowRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	err = bt.Remove(p.Sys + "." + p.Job)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db remove error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db remove error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobTimeWindowAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	m := new(module.MetaJobTimeWindowBean)
	m.Sys = p.Sys
	m.Job = p.Job
	m.Allow = p.Allow
	m.StartHour = p.StartHour
	m.EndHour = p.EndHour
	m.Description = p.Description
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Sys+"."+m.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobTimeWindowHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobTimeWindowBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_TIMEWINDOW)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	jobtw := bt.Get(p.Sys + "." + p.Job)
	if jobtw != nil {
		jtw := new(module.MetaJobTimeWindowBean)
		err := json.Unmarshal([]byte(jobtw.(string)), &jtw)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		if jtw.Enable != "1" {
			util.ApiResponse(response.ResponseWriter, 200, "", retlst)
			return
		}
		bhour := jtw.StartHour
		ehour := jtw.EndHour
		timeStrHour := int8(time.Now().Hour())
		if timeStrHour >= bhour || timeStrHour <= ehour {
			util.ApiResponse(response.ResponseWriter, 200, "", retlst)
			return
		} else {
			retlst = append(retlst, jtw)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobCmdGetAllHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdGetAllBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			mjc := new(module.MetaJobCmdBean)
			err := json.Unmarshal([]byte(v1.(string)), &mjc)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if mjc.Sys != jobpara.Sys || mjc.Job != jobpara.Job {
				continue
			}
			retlst = append(retlst, mjc)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobCmdGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Sys + "." + p.Job + "." + p.Step)
	if ib != nil {
		m := new(module.MetaJobCmdBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobCmdRemoveHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdRemoveBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	bt.Remove(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Step)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobCmdListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobCmdBean)
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

func (rrs *ResponseResource) FlowJobCmdAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobCmdAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	cb := new(module.MetaJobCmdBean)
	cb.Sys = p.Sys
	cb.Job = p.Job
	cb.Cmd = p.Cmd
	cb.Step = p.Step
	cb.Enable = p.Enable
	cb.Description = p.Description
	jsonstr, _ := json.Marshal(cb)

	err = bt.Set(p.Sys+"."+p.Job+"."+p.Step, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobCmdUpdateHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobCmdUpdateBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Step) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_CMD)
	defer bt.Close()

	fb0 := bt.Get(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Step)
	fb := new(module.MetaJobCmdBean)
	err = json.Unmarshal([]byte(fb0.(string)), &fb)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	fb.Type = jobpara.Type
	fb.Cmd = jobpara.Cmd
	fb.Step = jobpara.Step
	fb.Enable = jobpara.Enable
	jsonstr, _ := json.Marshal(fb)
	err = bt.Set(fb.Sys+"."+fb.Job+"."+fb.Step, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusUpdateSubmitHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdateSubmitBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	glog.Glog(LogF, fmt.Sprint(string(jsonstr)))
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusUpdateGoHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdateGoBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jobbn.SServer = jobpara.SServer
	jobbn.Sip = jobpara.Ip
	jobbn.Sport = jobpara.Port
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusUpdatePendingHandler(request *restful.Request, response *restful.Response) {
	rrs.Lock()
	defer rrs.Unlock()
	jobpara := new(module.MetaParaFlowJobStatusUpdatePendingBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	jobbn.Status = jobpara.Status
	jobbn.Sys = jobpara.Sys
	jobbn.Job = jobpara.Job
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusUpdateEndHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusUpdateEndBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	jobbn.EndTime = timeStr
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	jobpool.Remove(jobpara.FlowId + "," + jobpara.Sys + "," + jobpara.Job)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobStatusUpdateStartHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobStatusUpdateStartBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB)
	defer bt.Close()

	jobbn := new(module.MetaJobBean)
	jobn := bt.Get(jobpara.Sys + "." + jobpara.Job)
	if jobn != nil {
		err := json.Unmarshal([]byte(jobn.(string)), &jobbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("%v", err), nil)
			return
		}
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	jobbn.StartTime = timeStr
	jobbn.EndTime = ""
	jobbn.Status = jobpara.Status
	jsonstr, _ := json.Marshal(jobbn)
	err = bt.Set(jobbn.Sys+"."+jobbn.Job, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	mjm := new(module.MetaJobMemBean)
	mjm.FlowId = jobpara.FlowId
	mjm.Sys = jobbn.Sys
	mjm.Job = jobbn.Job
	mjm.CreateTime = timeStr
	mjm.Enable = "1"
	jobpool.Add(mjm.FlowId+","+mjm.Sys+","+mjm.Job, mjm)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobParameterRemoveHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterRemoveBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	bt.Remove(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Key)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobParameterGetHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterGetBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Sys != jobpara.Sys || m.Job != jobpara.Job || m.Key != jobpara.Key {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobParameterGetAllHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterGetAllBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	f, err := rrs.flowInfo(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow info %v error.%v", jobpara.FlowId, err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow info %v error.%v", jobpara.FlowId, err), nil)
		return
	}
	rrs.Lock()
	defer rrs.Unlock()
	retmap := make(map[string]interface{})
	//system
	bt0 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_PARAMETER)
	defer bt0.Close()
	strlist0 := bt0.Scan()
	bt0.Close()
	for _, v := range strlist0 {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			glog.Glog(LogF, fmt.Sprint(m))
			retmap[m.Key] = *m
		}
	}
	//flow
	bt1 := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_FLOW_PARAMETER)
	defer bt1.Close()
	strlist1 := bt1.Scan()
	bt1.Close()
	for _, v := range strlist1 {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			glog.Glog(LogF, fmt.Sprint(m))
			retmap[m.Key] = *m
		}
	}
	p := new(module.MetaJobParameterBean)
	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW_CTS
	p.Val = f.CreateTime
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW_CTS] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW_RCT
	p.Val = f.RunContext
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW_RCT] = *p

	//job
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobParameterBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get parameter failed.%v", err), nil)
				return
			}
			if m.Sys != jobpara.Sys || m.Job != jobpara.Job {
				continue
			}
			retmap[m.Key] = *m
		}
	}
	//force
	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_SYS
	p.Val = jobpara.Sys
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_SYS] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_JOB
	p.Val = jobpara.Job
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_JOB] = *p

	p.Sys = jobpara.Sys
	p.Job = jobpara.Job
	p.Key = util.CONST_FLOW
	p.Val = jobpara.FlowId
	p.Enable = util.CONST_ENABLE
	retmap[util.CONST_FLOW] = *p

	retlst := make([]interface{}, 0)
	for _, v := range retmap {
		retlst = append(retlst, v)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobParameterListHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaFlowJobParameterListBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
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

func (rrs *ResponseResource) FlowJobParameterAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobParameterAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 || len(p.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	m := new(module.MetaJobParameterBean)
	m.Sys = p.Sys
	m.Job = p.Job
	m.Key = p.Key
	m.Val = p.Val
	m.Description = p.Description
	m.Enable = p.Enable

	m.Type = "J"
	jsonstr, _ := json.Marshal(m)

	err = bt.Set(p.Sys+"."+p.Job+"."+p.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobParameterUpdateHandler(request *restful.Request, response *restful.Response) {
	jobpara := new(module.MetaParaBean)
	err := request.ReadEntity(&jobpara)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(jobpara.FlowId) == 0 || len(jobpara.Sys) == 0 || len(jobpara.Job) == 0 || len(jobpara.Key) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(jobpara.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_PARAMETER)
	defer bt.Close()

	fb0 := bt.Get(jobpara.Sys + "." + jobpara.Job + "." + jobpara.Key)
	m := new(module.MetaJobParameterBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	m.Val = jobpara.Val
	m.Description = jobpara.Description
	m.Enable = jobpara.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Sys+"."+m.Job+"."+m.Key, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) MstFlowRoutineJobRunningHeartListHandler(request *restful.Request, response *restful.Response) {

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

func (rrs *ResponseResource) MstFlowRoutineJobRunningHeartGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineJobRunningHeartRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) MstFlowRoutineJobRunningHeartAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) WorkerRoutineJobRunningHeartListHandler(request *restful.Request, response *restful.Response) {

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

func (rrs *ResponseResource) WorkerRoutineJobRunningHeartGetHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) WorkerRoutineJobRunningHeartRemoveHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) WorkerRoutineJobRunningHeartAddHandler(request *restful.Request, response *restful.Response) {
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

func (rrs *ResponseResource) FlowJobLogListHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogListBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.FlowId) == 0 || len(p.Sys) == 0 || len(p.Job) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	strlist := bt.Scan()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaJobLogBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			if m.Sys != p.Sys || m.Job != p.Job {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobLogGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	ib := bt.Get(p.Id)
	if ib != nil {
		m := new(module.MetaJobLogBean)
		err := json.Unmarshal([]byte(ib.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
		}
		retlst = append(retlst, m)
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResource) FlowJobLogRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
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

func (rrs *ResponseResource) FlowJobLogAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	m := new(module.MetaJobLogBean)
	m.Id = p.Id
	m.Sys = p.Sys
	m.Job = p.Job
	m.StartTime = p.StartTime
	m.EndTime = p.EndTime
	m.SServer = p.SServer
	m.Sip = p.Sip
	m.Sport = p.Sport
	m.Step = p.Step
	m.Content = make([]string, 0)
	m.ExitCode = p.ExitCode
	m.Cmd = p.Cmd
	for _, v := range p.Content {
		m.Content = append(m.Content, v)
	}
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

func (rrs *ResponseResource) FlowJobLogAppendHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaFlowJobLogAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 || len(p.FlowId) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}
	dbf, err := rrs.flowDbFile(p.FlowId)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("get flow db file error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get flow db file error.%v", err), nil)
		return
	}
	bt := db.NewBoltDB(dbf, util.TABLE_AUTO_JOB_LOG)
	defer bt.Close()

	b := bt.Get(p.Id)
	m := new(module.MetaJobLogBean)
	if b != nil {
		err := json.Unmarshal([]byte(b.(string)), &m)
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("get in db error.%v", err))
			util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get in db error.%v", err), nil)
			return
		}
	} else {
		m.Id = p.Id
		m.Sys = p.Sys
		m.Job = p.Job
		m.StartTime = p.StartTime
		m.SServer = p.SServer
		m.Sip = p.Sip
		m.Sport = p.Sport
		m.Step = p.Step
		m.Cmd = p.Cmd
		m.Content = make([]string, 0)
	}
	for _, v := range p.Content {
		m.Content = append(m.Content, v)
	}
	m.EndTime = p.EndTime
	m.ExitCode = p.ExitCode
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

var (
	conf             *module.MetaApiServerBean
	jobpool          = db.NewMemDB()
	slvMap           = db.NewMemDB()
	jobspool         = db.NewMemDB()
	flowserverspool  = db.NewMemDB()
	ringPendingSpool = db.NewMemDB()
	ringGoSpool      = db.NewMemDB()
	LogF             string
	//mstMap  map[string]interface{}
	mstMap                = db.NewMemDB()
	statusPendingOffset   int
	statusGoOffset        int
	statusPendingHashRing *util.Consistent
	statusGoHashRing      *util.Consistent
	rru                   *ResponseResourceUser
)

func NewApiServer(cfg string) {
	conf = new(module.MetaApiServerBean)
	yamlFile, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Printf("error: %s", err)
		return
	}
	err = yaml.UnmarshalStrict(yamlFile, conf)
	if err != nil {
		log.Printf("error: %s", err)
		return
	}
	LogF = conf.HomeDir + "/handle_${" + util.ENV_VAR_DATE + "}.log"
	statusPendingHashRing = util.NewConsistent()
	for i := 0; i < 6; i++ {
		si := fmt.Sprintf("%d", i)
		statusPendingHashRing.Add(util.NewNode(i, si, 1))
	}
	statusGoHashRing = util.NewConsistent()
	for i := 0; i < 6; i++ {
		si := fmt.Sprintf("%d", i)
		statusGoHashRing.Add(util.NewNode(i, si, 1))
	}
	statusPendingOffset = 0
	statusGoOffset = 0
	go func() {
		JobPool()
	}()
	rru = NewResponseResourceUser()
	HttpServer()
}

func JobPool() {
	for {
		glog.Glog(LogF, fmt.Sprintf("Check Job Pool..."))
		url := fmt.Sprintf("http://%v:%v/api/v1/job/pool/ls?accesstoken=%v", conf.ApiServerIp, conf.ApiServerPort, conf.AccessToken)
		jsonstr, err := util.Api_RequestPost(url, "{}")
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		retbn := new(module.RetBean)
		err = json.Unmarshal([]byte(jsonstr), &retbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		if retbn.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		if retbn.Data == nil {
			glog.Glog(LogF, fmt.Sprint("get pending status job err."))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		retarr := (retbn.Data).([]interface{})
		serverlst := ObtSlvRunningJobCnt()
		for i := 0; i < len(retarr); i++ {
			v := retarr[i].(map[string]interface{})
			s := ObtJobServer(serverlst, v["server"].(string))
			if len(s) < 1 {
				glog.Glog(LogF, fmt.Sprintf("%v, %v, %v job no free server running,wait for next time.", v["flowid"], v["sys"], v["job"]))
				time.Sleep(time.Duration(10) * time.Second)
				continue
			}
			SubmitJob(v, s[0])
		}
		if len(retarr) == 0 {
			glog.Glog(LogF, fmt.Sprint("job pool empty ,no running job ."))
		} else {
			glog.Glog(LogF, fmt.Sprintf("do with job cnt %v", len(retarr)))
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func SubmitJob(jobstr interface{}, slvstr interface{}) bool {
	s := slvstr.(map[string]interface{})
	jp := jobstr.(map[string]interface{})
	glog.Glog(LogF, fmt.Sprintf("%v update %v %v job status %v.", jp["flowid"], jp["sys"], jp["job"], util.STATUS_AUTO_GO))
	url0 := fmt.Sprintf("http://%v:%v/api/v1/flow/job/status/update/go?accesstoken=%v", conf.ApiServerIp, conf.ApiServerPort, conf.AccessToken)
	m := new(module.MetaParaFlowJobStatusUpdateGoBean)
	m.FlowId = jp["flowid"].(string)
	m.Sys = jp["sys"].(string)
	m.Job = jp["job"].(string)
	m.Status = util.STATUS_AUTO_GO
	m.SServer = s["workerid"].(string)
	m.Ip = s["ip"].(string)
	m.Port = s["port"].(string)
	jsonstr0, err := json.Marshal(m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	jsonstr, err := util.Api_RequestPost(url0, string(jsonstr0))
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url %v return status code:%v,msg: %v", url0, retbn.Status_Code, retbn.Status_Txt))
		return false
	}

	glog.Glog(LogF, fmt.Sprintf("rm from job pool %v,%v,%v.", jp["sys"], jp["job"], jp["flowid"]))
	url := fmt.Sprintf("http://%v:%v/api/v1/job/pool/rm?accesstoken=%v", conf.ApiServerIp, conf.ApiServerPort, conf.AccessToken)
	para := fmt.Sprintf("{\"sys\":\"%v\",\"job\":\"%v\",\"flowid\":\"%v\"}", jp["sys"], jp["job"], jp["flowid"])
	jsonstr1, err := util.Api_RequestPost(url, para)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	retbn1 := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr1), &retbn1)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return false
	}
	if retbn1.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url %v return status code:%v", url, retbn1.Status_Code))
		return false
	}

	return true
}

func ObtSlvRunningJobCnt() []interface{} {
	url := fmt.Sprintf("http://%v:%v/api/v1/worker/cnt/exec?accesstoken=%v", conf.ApiServerIp, conf.ApiServerPort, conf.AccessToken)
	jsonstr, err := util.Api_RequestPost(url, "{}")
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return nil
	}
	retbn := new(module.RetBean)
	err = json.Unmarshal([]byte(jsonstr), &retbn)
	if err != nil {
		glog.Glog(LogF, fmt.Sprint(err))
		return nil
	}
	if retbn.Status_Code != 200 {
		glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
		return nil
	}
	if retbn.Data == nil {
		glog.Glog(LogF, fmt.Sprint("get pending status job err."))
		return nil
	}
	retarr := (retbn.Data).([]interface{})
	return retarr
}

func ObtJobServer(arr []interface{}, slvid string) []interface{} {
	retlst := make([]interface{}, 0)
	for i := 0; i < len(arr); i++ {
		v := arr[i].(map[string]interface{})
		maxcnt, err := strconv.Atoi(v["maxcnt"].(string))
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("string conv int %v err.%v", v["maxcnt"], err))
			continue
		}
		runningcnt, err := strconv.Atoi(v["runningcnt"].(string))
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("string conv int %v err.%v", v["runningcnt"], err))
			continue
		}
		currentcnt, err := strconv.Atoi(v["currentcnt"].(string))
		if err != nil {
			glog.Glog(LogF, fmt.Sprintf("string conv int %v err.%v", v["currentcnt"], err))
			continue
		}
		if len(slvid) > 0 {
			if v["slaveid"] == slvid {
				retlst = make([]interface{}, 0)
				retlst = append(retlst, arr[i])
				break
			}
		} else if 0 < (maxcnt - runningcnt - currentcnt) {
			retlst = append(retlst, arr[i])
		}
	}
	return retlst
}

func ServerRoutineStatus() {
	for {
		url := fmt.Sprintf("http://%v:%v/api/v1/mst/heart/ls?accesstoken=%v", conf.ApiServerIp, conf.ApiServerPort, conf.AccessToken)
		jsonstr, err := util.Api_RequestPost(url, "{}")
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		retbn := new(module.RetBean)
		err = json.Unmarshal([]byte(jsonstr), &retbn)
		if err != nil {
			glog.Glog(LogF, fmt.Sprint(err))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		if retbn.Status_Code != 200 {
			glog.Glog(LogF, fmt.Sprintf("post url return status code:%v", retbn.Status_Code))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		if retbn.Data == nil {
			glog.Glog(LogF, fmt.Sprint("get pending status job err."))
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		mharr := (retbn.Data).([]module.MetaMstHeartBean)
		for _, mh := range mharr {
			// 建立连接到gRPC服务
			conn, err := grpc.Dial(mh.Ip+":"+mh.Port, grpc.WithInsecure())
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("did not connect: %v", err))
				continue
			}
			// 函数结束时关闭连接
			defer conn.Close()

			// 创建Waiter服务的客户端
			t := gproto.NewFlowMasterClient(conn)

			mifb := new(module.MetaInstanceFlowBean)
			jsonstr, err := json.Marshal(mifb)
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("json marshal %v", err))
				continue
			}
			// 调用gRPC接口
			tr, err := t.FlowStatus(context.Background(), &gproto.Req{JsonStr: string(jsonstr)})
			if err != nil {
				glog.Glog(LogF, fmt.Sprintf("could not greet: %v", err))
				continue
			}
			//change job status
			if tr.Status_Code != 200 {
				glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
				continue
			}
			if len(tr.Data) == 0 {
				glog.Glog(LogF, fmt.Sprint(tr.Status_Txt))
				continue
			}
			mf := new([]module.MetaFlowStatusBean)
			err = json.Unmarshal([]byte(tr.Data), &mf)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
			}
			for _, v := range *mf {
				glog.Glog(LogF, fmt.Sprint(v))
			}
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func HttpServer() {
	restful.Filter(globalOauth)

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Access-Control-Allow-*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	//rrs := ResponseResource{Data: make(map[string]interface{})}
	rrs := ResponseResource{}
	restful.DefaultContainer.Add(rrs.WebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/home/k8s/Go/ws/src/github.com/lhzd863/autoflow/swagger-ui-3.22.0/dist"))))

	log.Printf("Get the API using http://localhost:" + conf.Port + "/apidocs.json")
	log.Printf("Open Swagger UI using http://" + conf.ApiServerIp + ":" + conf.Port + "/apidocs/?url=http://localhost:" + conf.Port + "/apidocs.json")
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Autoflow",
			Description: "Resource for managing Api",
			Contact: &spec.ContactInfo{
				Name:  "lhzd863",
				Email: "lhzd863@126.com",
				URL:   "http://lhzd863.com",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://lhzd863.com",
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "image",
		Description: "Managing image"}}}
}

// Global Filter
func globalOauth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	u, err := url.Parse(req.Request.URL.String())
	if err != nil {
		log.Println("parse url error")
		return
	}
	if u.Path == "/api/v1/login" || u.Path == "/api/v1/register" || u.Path == "/apidocs" || u.Path == "/apidocs.json" {
		chain.ProcessFilter(req, resp)
		return
	}
	reqParams, err := url.ParseQuery(req.Request.URL.RawQuery)
	if err != nil {
		log.Printf("Failed to decode: %v", err)
		util.ApiResponse(resp.ResponseWriter, 700, fmt.Sprintf("Failed decode error.%v", err), nil)
		return
	}
	tokenstring := fmt.Sprint(reqParams["accesstoken"][0])
	var claimsDecoded map[string]interface{}
	decodeErr := jwt.Decode([]byte(tokenstring), &claimsDecoded, []byte(conf.JwtKey))
	if decodeErr != nil {
		log.Printf("Failed to decode: %s (%s)", decodeErr, tokenstring)
		util.ApiResponse(resp.ResponseWriter, 700, fmt.Sprintf("Failed to decode.%v", decodeErr), nil)
		return
	}

	exp := claimsDecoded["exp"].(float64)
	exp1, _ := strconv.ParseFloat(fmt.Sprintf("%v", time.Now().Unix()+0), 64)

	if (exp - exp1) < 0 {
		log.Printf("Failed to decode: %v %v %v", exp, exp1, (exp - exp1))
		util.ApiResponse(resp.ResponseWriter, 700, "Not Authorized AccessToken Expired ,Please login", nil)
		return
	}
	//fmt.Println((exp - exp1))
	chain.ProcessFilter(req, resp)
}
