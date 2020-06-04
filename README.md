# autoflow

## 简介
```
分布调度系统服务端
```
## 依赖
```
-github.com/satori/go.uuid
-google.golang.org/grpc
-github.com/emicklei/go-restful
-go.etcd.io/bbolt
-jwt
-workpool

```
## 名词
```
-镜像: 带有存储作业配置，依赖，参数，触发等相关配置的bbolt存储文件，建立目录创建ID和打上标记形成对象
-实例: 镜像文件被具体到某个批量，存储文件从镜像文件copy，不会影响镜像配置文件
-作业: 实现例具体执行操作对象

```
## 功能
```
-API
 -api文档采用swagger
 -实例作业的添加，删除，修改
 -镜像文件库添加，删除，修改
 -实例创建，启动，停止
 
-Mgr
 -检查Pending状态作业
 -检查Go状态作业
 
-Worker
 -执行作业
 -作业执行LOG信息返回
 -参数替换和环境变量设置顺序:系统->实例->作业
 
```

## Online Demo

[在线 Demo](https://122.51.161.53:12300)
