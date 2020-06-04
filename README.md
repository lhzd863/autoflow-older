# autoflow

## 简介

autoflow数据批量分布式作业调度系,支持批量调度之间无干扰，元数据相互对立，使用元数据能够快速恢复一个相同运行环境，系统间通信采用[go-restful](https://github.com/emicklei/go-restful)进行交互，数据采用json方式进行存储，节点和节点之间通过grpc通信，消息格式采用protobuffer。系统采用服务端和客户端相分离方式，只有同时使用服务端和客户端才能完成整套调度系统搭建。服务端元数据存储采用[bbolt](https://github.com/etcd-io/bbolt)无需安装，配置相关参数可以直接使用，客户端在[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin) 基础上开发完成。

## 依赖

- [go.uuid](https://github.com/satori/Go.uuid) 唯一ID生成
- google.golang.org/grpc
- [go-restful](https://github.com/emicklei/go-restful) restful web框架
- [bbolt](https://github.com/etcd-io/bbolt) 数据存储
- [jwt](https://github.com/robbert229/jwt) token安全验证
- [workpool](https://github.com/goinggo/workpool) 线程池


## 名词
```
-镜像: 带有存储作业配置，依赖，参数，触发等相关配置的bbolt存储文件，建立目录创建ID和打上标记形成对象
-实例: 镜像文件被具体到某个批量，存储文件从镜像文件copy，不会影响镜像配置文件
-作业: 实现例具体执行操作对象

```
## 功能
```
- API
 - api文档采用swagger
 - 实例作业的添加，删除，修改
 - 镜像文件库添加，删除，修改
 - 实例创建，启动，停止
 - 实例线程增加和减少
  - 指定节点上增加，减少指定数量线程
  - 增加减少具体线程名称
 
- Mgr
 - 检查Pending状态作业
 - 检查Go状态作业
 
- Worker
 - 执行作业
 - 作业执行LOG信息返回
 - 参数替换和环境变量设置顺序:系统->实例->作业
 
```

## Online Demo

[在线 Demo](https://122.51.161.53:12300)

## 项目声明：
本软件只供学习交流使用，勿作为商业用途,如有任何问题联系lhzd863@126.com
