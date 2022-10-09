# websocket #

## 注意信息
### 健康检查前缀/health,硬编码兼容修改
1. 本次修改是为了适配k8s clb的http 健康检查
2. clb检查频繁, 过滤health path是为了减少健康检查的日志信息

## 模块介绍

websocket初始使用http发送握手信息，握手成功后接管底层的tcp连接，发送/接收按规定格式封装的帧。

websocket 模块基于gobwas库进行封装。实现事件驱动方式接口。

## 实现架构

![](../../docs/images/websocket-layer.png)

## 事件处理流程

![](../../docs/images/websocket-loop.png)

## cmd模式生命周期

![](../../docs/images/websocket-lifecircle.png)

## 可自定义属性

- 分片设置

模块提供自动分片写的机制，默认分片大小为1024，可通过设置FrameSize的值自定义分片大小。

- 心跳设置


模块启用ping帧监测心跳，通过记录收到客户端回复的pong帧时间，判断连接是否超时。心跳发送间隔默认为3秒， pong帧超时事件默认是3次心跳发送间隔，可通过设置KeepAliveInterval设置心跳发送间隔，可通过设置LastPongTimeout设置pong帧超时时间。





