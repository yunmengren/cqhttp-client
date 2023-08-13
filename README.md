# cqhttp-client
基于cqhttp的群消息监控及定时任务

## 功能
- 设置关键词，监控账号内所有接收到的消息，包含关键词的消息发送到指定账号
- 创建定时任务，定时向指定账号发送消息

## 安装部署
1. 部署签名服务器https://hub.docker.com/r/xzhouqd/qsign
2. 部署cqhttp服务https://docs.go-cqhttp.org/guide/quick_start.html#%E4%BD%BF%E7%94%A8
3. 配置端口连接，连接服务

## 问题
- 定时发送消息，容易被系统检测为异常行为，导致封号