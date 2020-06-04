# autoflow
## 依赖
```
-github.com/satori/go.uuid
-google.golang.org/grpc
-github.com/emicklei/go-restful
-go.etcd.io/bbolt
-jwt
-workpool

```

## 功能
```
-API
 -api文档采用swagger
 -实例作业的添加，删除，修改
 
-Mgr
 -检查Pending状态作业
 -检查Go状态作业
 
-Worker
 -执行作业
 -作业执行LOG信息返回
 -参数替换和环境变量设置顺序:系统->实例->作业
 
```
