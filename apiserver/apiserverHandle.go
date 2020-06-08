package apiserver

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	ws.Route(ws.POST("/system/user/info").To(rru.SystemUserInfoHandler).
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
	ws.Route(ws.POST("/system/parameter/add").To(rrt.SystemParameterAddHandler).
		// docs
		Doc("新增系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/rm").To(rrt.SystemParameterRemoveHandler).
		// docs
		Doc("删除系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/get").To(rrt.SystemParameterGetHandler).
		// docs
		Doc("获取系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/ls").To(rrt.SystemParameterListHandler).
		// docs
		Doc("列表系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/parameter/update").To(rrt.SystemParameterUpdateHandler).
		// docs
		Doc("更新系统参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system"}
	ws.Route(ws.POST("/sys/lsport").To(rrt.SysListPortHandler).
		// docs
		Doc("list port").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"image"}
	ws.Route(ws.POST("/image/add").To(rri.ImageAddHandler).
		// docs
		Doc("新增镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/rm").To(rri.ImageRemoveHandler).
		// docs
		Doc("删除镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/ls").To(rri.ImageListHandler).
		// docs
		Doc("列表镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/get").To(rri.ImageGetHandler).
		// docs
		Doc("获取镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/image/update").To(rri.ImageUpdateHandler).
		// docs
		Doc("更新镜像文件").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaImageUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"instance"}
	ws.Route(ws.POST("/instance/create").To(rrf.InstanceCreateHandler).
		// docs
		Doc("实例创建").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaJobFlowBean{ImageId: "", ProcessNum: ""}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/start").To(rrf.InstanceStartHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/stop").To(rrf.InstanceStopHandler).
		// docs
		Doc("实例停止").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/ls").To(rrf.InstanceListHandler).
		// docs
		Doc("实例列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/ls/status").To(rrf.InstanceListStatusHandler).
		// docs
		Doc("实例状态列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/instance/rm").To(rrf.InstanceRemoveHandler).
		// docs
		Doc("实例列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow"}
	ws.Route(ws.POST("/mst/instance/start").To(rrf.MstInstanceStartHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowStartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/instance/stop").To(rrf.MstInstanceStopHandler).
		// docs
		Doc("实例开启").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowStopBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-routine"}
	ws.Route(ws.POST("/flow/routine/add").To(rrf.FlowRoutineAddHandler).
		// docs
		Doc("实例新增处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowRoutineAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/routine/sub").To(rrf.FlowRoutineSubHandler).
		// docs
		Doc("实例线程删除").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowRoutineSubBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/routine/status").To(rrf.FlowRoutineStatusHandler).
		// docs
		Doc("流获取").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine"}
	ws.Route(ws.POST("/mst/flow/routine/start").To(rrf.MstFlowRoutineStartHandler).
		// docs
		Doc("指定mst实例新增处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/stop").To(rrf.MstFlowRoutineStopHandler).
		// docs
		Doc("指定mst实例停止处理线程").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow"}
	ws.Route(ws.POST("/flow/ls").To(rrf.FlowListHandler).
		// docs
		Doc("列表实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/get").To(rrf.FlowGetHandler).
		// docs
		Doc("获取实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/update").To(rrf.FlowUpdateHandler).
		// docs
		Doc("更新实例").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/status/update").To(rrf.FlowStatusUpdateHandler).
		// docs
		Doc("更新实例状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowStatusUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst"}
	ws.Route(ws.POST("/mst/flow/rm").To(rrf.MstFlowRemoveHandler).
		// docs
		Doc("删除实例MST映射").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-parameter"}
	ws.Route(ws.POST("/flow/parameter/add").To(rrf.FlowParameterAddHandler).
		// docs
		Doc("新增实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/ls").To(rrf.FlowParameterListHandler).
		// docs
		Doc("列表实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/get").To(rrf.FlowParameterGetHandler).
		// docs
		Doc("获取实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/rm").To(rrf.FlowParameterRemoveHandler).
		// docs
		Doc("删除实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/parameter/update").To(rrf.FlowParameterUpdateHandler).
		// docs
		Doc("更新实例参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-pool"}
	ws.Route(ws.POST("/job/pool/add").To(rrt.JobPoolAddHandler).
		// docs
		Doc("新增作业到系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/rm").To(rrt.JobPoolRemoveHandler).
		// docs
		Doc("删除作业从系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/get").To(rrt.JobPoolGetHandler).
		// docs
		Doc("获取作业从系统作业池").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/job/pool/ls").To(rrt.JobPoolListHandler).
		// docs
		Doc("列表在系统作业池中作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker"}
	ws.Route(ws.POST("/worker/heart/add").To(rrw.WorkerHeartAddHandler).
		// docs
		Doc("注册执行节点").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/rm").To(rrw.WorkerHeartRemoveHandler).
		// docs
		Doc("节点删除").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/ls").To(rrw.WorkerHeartListHandler).
		// docs
		Doc("节点列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/heart/get").To(rrw.WorkerHeartGetHandler).
		// docs
		Doc("节点获取").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaWorkerHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker-cnt"}
	ws.Route(ws.POST("/worker/cnt/add").To(rrw.WorkerCntAddHandler).
		// docs
		Doc("节点作业数增加").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaWorkerHeartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/cnt/exec").To(rrw.WorkerExecCntHandler).
		// docs
		Doc("节点执行作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job"}
	ws.Route(ws.POST("/flow/job/ls").To(rrj.FlowJobListHandle).
		// docs
		Doc("列表实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/get").To(rrj.FlowJobGetHandle).
		// docs
		Doc("获取实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/add").To(rrj.FlowJobAddHandle).
		// docs
		Doc("新增实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/update").To(rrj.FlowJobUpdateHandle).
		// docs
		Doc("更新实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/rm").To(rrj.FlowJobRemoveHandle).
		// docs
		Doc("删除实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-status-get"}
	ws.Route(ws.POST("/flow/job/status/get/pending").To(rrj.FlowJobStatusGetPendingHandle).
		// docs
		Doc("获取Pending状态实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusGetPendingBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/get/go").To(rrj.FlowJobStatusGetGoHandle).
		// docs
		Doc("获取Go状态实例作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusGetGoBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-heart"}
	ws.Route(ws.POST("/mst/heart/add").To(rrl.MstHeartAddHandler).
		// docs
		Doc("新增Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/ls").To(rrl.MstHeartListHandler).
		// docs
		Doc("列表Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/rm").To(rrl.MstHeartRemoveHandler).
		// docs
		Doc("删除Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/heart/get").To(rrl.MstHeartGetHandler).
		// docs
		Doc("获取Mst心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine-heart"}
	ws.Route(ws.POST("/mst/flow/routine/heart/add").To(rrl.MstFlowRoutineHeartAddHandler).
		// docs
		Doc("新增Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/ls").To(rrl.MstFlowRoutineHeartListHandler).
		// docs
		Doc("列表Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/rm").To(rrl.MstFlowRoutineHeartRemoveHandler).
		// docs
		Doc("删除Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/heart/get").To(rrl.MstFlowRoutineHeartGetHandler).
		// docs
		Doc("获取Mst节点实例线程心跳信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaMstFlowRoutineHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-dependency"}
	ws.Route(ws.POST("/flow/job/dependency").To(rrj.FlowJobDependencyHandler).
		// docs
		Doc("流作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/ls").To(rrj.FlowJobDependencyListHandler).
		// docs
		Doc("列表作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/rm").To(rrj.FlowJobDependencyRemoveHandler).
		// docs
		Doc("删除作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/get").To(rrj.FlowJobDependencyGetHandler).
		// docs
		Doc("获取作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/update").To(rrj.FlowJobDependencyUpdateHandler).
		// docs
		Doc("更新作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/dependency/add").To(rrj.FlowJobDependencyAddHandler).
		// docs
		Doc("新增作业依赖").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobDependencyAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-timewindow"}
	ws.Route(ws.POST("/flow/job/timewindow").To(rrj.FlowJobTimeWindowHandler).
		// docs
		Doc("实例作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/add").To(rrj.FlowJobTimeWindowAddHandler).
		// docs
		Doc("新增作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/rm").To(rrj.FlowJobTimeWindowRemoveHandler).
		// docs
		Doc("删除作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/ls").To(rrj.FlowJobTimeWindowListHandler).
		// docs
		Doc("列表作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/get").To(rrj.FlowJobTimeWindowGetHandler).
		// docs
		Doc("获取作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/timewindow/update").To(rrj.FlowJobTimeWindowUpdateHandler).
		// docs
		Doc("更新作业时间窗口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobTimeWindowUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-stream"}
	ws.Route(ws.POST("/flow/job/stream/job").To(rrj.FlowJobStreamJobHandler).
		// docs
		Doc("触发作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamJobBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/job/get").To(rrj.FlowJobStreamJobGetHandler).
		// docs
		Doc("触发作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamJobListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/add").To(rrj.FlowJobStreamAddHandler).
		// docs
		Doc("新增作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/rm").To(rrj.FlowJobStreamRemoveHandler).
		// docs
		Doc("删除作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/ls").To(rrj.FlowJobStreamListHandler).
		// docs
		Doc("列表作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/get").To(rrj.FlowJobStreamGetHandler).
		// docs
		Doc("获取作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/stream/update").To(rrj.FlowJobStreamUpdateHandler).
		// docs
		Doc("更新作业触发").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStreamUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-status-update"}
	ws.Route(ws.POST("/flow/job/status/update/submit").To(rrj.FlowJobStatusUpdateSubmitHandler).
		// docs
		Doc("更新作业状态为Submit").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateSubmitBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/go").To(rrj.FlowJobStatusUpdateGoHandler).
		// docs
		Doc("更新作业状态为Go").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateGoBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/pending").To(rrj.FlowJobStatusUpdatePendingHandler).
		// docs
		Doc("更新作业状态为Pending").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdatePendingBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/start").To(rrj.FlowJobStatusUpdateStartHandler).
		// docs
		Doc("更新作业开始运行相关信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateStartBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/status/update/end").To(rrj.FlowJobStatusUpdateEndHandler).
		// docs
		Doc("更新作业结束运行相关信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobStatusUpdateEndBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-cmd"}
	ws.Route(ws.POST("/flow/job/cmd/add").To(rrj.FlowJobCmdAddHandler).
		// docs
		Doc("新增作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/getall").To(rrj.FlowJobCmdGetAllHandler).
		// docs
		Doc("获取作业所有脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Reads(module.MetaParaFlowJobCmdGetAllBean{}).
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/get").To(rrj.FlowJobCmdGetHandler).
		// docs
		Doc("获取作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Reads(module.MetaParaFlowJobCmdGetBean{}).
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/ls").To(rrj.FlowJobCmdListHandler).
		// docs
		Doc("列表作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/update").To(rrj.FlowJobCmdUpdateHandler).
		// docs
		Doc("更新作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/cmd/rm").To(rrj.FlowJobCmdRemoveHandler).
		// docs
		Doc("删除作业脚本").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobCmdRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-parameter"}
	ws.Route(ws.POST("/flow/job/parameter/add").To(rrj.FlowJobParameterAddHandler).
		// docs
		Doc("新增作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/get").To(rrj.FlowJobParameterGetHandler).
		// docs
		Doc("获取作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/getall").To(rrj.FlowJobParameterGetAllHandler).
		// docs
		Doc("列表系统实例作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterGetAllBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/ls").To(rrj.FlowJobParameterListHandler).
		// docs
		Doc("列表作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/update").To(rrj.FlowJobParameterUpdateHandler).
		// docs
		Doc("更新作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterUpdateBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/parameter/rm").To(rrj.FlowJobParameterRemoveHandler).
		// docs
		Doc("删除作业参数").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobParameterRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"flow-job-log"}
	ws.Route(ws.POST("/flow/job/log/add").To(rrj.FlowJobLogAddHandler).
		// docs
		Doc("新增作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/get").To(rrj.FlowJobLogGetHandler).
		// docs
		Doc("获取作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/ls").To(rrj.FlowJobLogListHandler).
		// docs
		Doc("列表作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogListBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/rm").To(rrj.FlowJobLogRemoveHandler).
		// docs
		Doc("删除作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/flow/job/log/append").To(rrj.FlowJobLogAppendHandler).
		// docs
		Doc("追加作业日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaFlowJobLogAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"mst-flow-routine-job-running-heart"}
	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/add").To(rrl.MstFlowRoutineJobRunningHeartAddHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/rm").To(rrl.MstFlowRoutineJobRunningHeartRemoveHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/get").To(rrl.MstFlowRoutineJobRunningHeartGetHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/mst/flow/routine/job/running/heart/ls").To(rrl.MstFlowRoutineJobRunningHeartListHandler).
		// docs
		Doc("新增Running作业").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"worker-routine-job-running-heart"}
	ws.Route(ws.POST("/worker/routine/job/running/heart/ls").To(rrw.WorkerRoutineJobRunningHeartListHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/get").To(rrw.WorkerRoutineJobRunningHeartGetHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartGetBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/rm").To(rrw.WorkerRoutineJobRunningHeartRemoveHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/worker/routine/job/running/heart/add").To(rrw.WorkerRoutineJobRunningHeartAddHandler).
		// docs
		Doc("Worker Running作业列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemMstFlowRoutineJobRunningHeartAddBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	tags = []string{"system-ring"}
	ws.Route(ws.POST("/system/ring/pending/ls").To(rrt.SystemRingPendingListHandler).
		// docs
		Doc("列表Pending Ring").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/go/ls").To(rrt.SystemRingGoListHandler).
		// docs
		Doc("列表Go Ring").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/pending/rm").To(rrt.SystemRingPendingRemoveHandler).
		// docs
		Doc("删除Pending Ring信息").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.QueryParameter("accesstoken", "access token").DataType("string")).
		Reads(module.MetaParaSystemRingPendingRemoveBean{}).
		Writes(module.RetBean{}). // on the response
		Returns(200, "OK", module.RetBean{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.POST("/system/ring/go/rm").To(rrt.SystemRingGoRemoveHandler).
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
	rri                   *ResponseResourceImage
	rrf                   *ResponseResourceFlow
	rrt                   *ResponseResourceSystem
	rrl                   *ResponseResourceLeader
	rrw                   *ResponseResourceWorker
	rrj                   *ResponseResourceJob
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
	rri = NewResponseResourceImage()
	rrf = NewResponseResourceFlow()
	rrt = NewResponseResourceSystem()
	rrl = NewResponseResourceLeader()
	rrw = NewResponseResourceWorker()
	rrj = NewResponseResourceJob()
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
