# cuckoo

基于 [bitcapybara/raft](https://github.com/bitcapybara/raft) 库实现的分布式任务调度系统，由于使用了 raft 算法使得任务调度完全基于内存进行，相比于传统基于数据库的调度系统进一步提高性能

## 一、架构
* 分为服务端和客户端，每一个服务端是一个 raft 节点
* 服务端依赖作为单独进程运行，进行任务管理和调度
* 客户端是执行任务的实体，通过定期向服务端发送心跳注册到服务端的任务执行组中
* 指定任务的路由策略后，由服务端来计算出任务由那个客户端执行

## 二、功能

### 调度类型
* `cron`：基于 [robfig/cron](https://github.com/robfig/cron) 实现，支持标准 cron 表达式及各种扩展写法
* `fixedRate`：可指定首次延迟和任务之间执行的延迟，间隔固定时长触发

### 路由类型
* `first`：任务执行组中的第一个节点
* `last`：任务执行组中的最后一个节点
* `random`：在任务执行组中随机选择一个节点
* `round`：轮询任务执行组中的节点

## 三、服务端需要实现的接口

### 1. bitcapbara/raft 库中的接口

#### Transport
> 在 raft 内部调用此接口的各个方法用于网络通信，比如发送心跳，日志复制，领导者选举，发送快照等。

#### RaftStatePersister
> 在 raft 内部调用此接口来持久化和加载内部状态数据，包括 term，votedFor及日志条目。

#### SnapshotPersister
> 在 raft 内部调用此接口来持久化和加载快照数据。

#### Logger
> 在 raft 内部调用此接口来打印日志。

### 2. bitcapbara/cuckoo 库中的接口

#### JobDispatcher
> 在 cuckoo 内部调用此接口向客户端发送任务执行请求

## 四、使用

### 服务端
* 依赖 `cuckoo/server`, 新建 `cuckoo.server.Server` 对象，代表当前节点
* 调用 `cuckoo.server..Start()` 启动服务端进程
* 在开放的 HTTP/RPC 接口中调用 `cuckoo.server.Server` 的公共方法

## 五、示例
[simpleSched](https://github.com/bitcapybara/simpleSched) 项目是此项目的一个简单实现，网络通信采用 `gin` 库