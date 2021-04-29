# cuckoo

基于 [bitcapybara/raft](https://github.com/bitcapybara/raft) 库实现的分布式任务调度系统，由于使用了 raft 算法使得任务调度完全基于内存进行，相比于传统基于数据库的调度系统进一步提高性能

### 架构
分为 server 端和 client 端，每一个 server 是一个 raft 节点
* 业务代码依赖 cuckoo/client，调用任务提交接口向服务端发送任务增删改查请求
* 服务端依赖 cuckoo/server 作为单独进程运行，接收客户端的请求进行任务管理和调度
* client 定期发送心跳来向 server 注册自身地址保活
* 任务到期后，server 通过向 client 发送请求来触发任务的执行，任务的实际执行者是 client