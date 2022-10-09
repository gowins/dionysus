
# healthy

## 模块介绍

healthy模块为ctl，wscmd， grpccmd提供健康检查方法。
### ctl健康检查

ctl健康检查原理为，每隔一段时间在指定目录下更新指定的文件，若文件更新时间超过一定时间间隔，则认为不健康，否则则为健康状态。

命令形式为`ctl healthx`

### wscmd健康检查

websocket server的健康检查机制是将websocket server的端口号写入指定文件，建立websocket客户端向服务器建立websocket连接，连接建立成功，则为健康状态，否则为非健康状态

命令形式为`websocket healthx`

### grpc健康检查

grpc server健康检查方式与websocket机制相同。通过客户端与服务器通信判断健康状态。

命令形式为`grpc healthx`

### 获取服务ADDR实现

websocket与grpc需要知道server ADDR，才能与服务端通信，所以这两种方式需要获取服务端口号。在服务启动的prerun中调用```healthy.WritePortFile(Addr, healthy.XXPortFile)```写入容器对应文件中。
执行healthy cmd时先调用SetPort从PortFile读取端口后设置对应ADDR。然后进行服务访问检查服务健康状况CheckHealthy。


## 使用方法

启用ctl， wscmd, grpccmd模块，healthx子命令默认可用。

### k8s健康检查
实际上这个健康检查机制主要为了应用部署在k8s中为容器的健康检查提供方法。实际使用中，yml文件可如下配置：

#### ctl 程序
```yaml
livenessProbe:
      exec:
        command:
        - /main
        - ctl
        - healthx
      initialDelaySeconds: 5
      periodSeconds: 5
```

#### websocket程序
```yaml
livenessProbe:
      exec:
        command:
        - /main
        - websocket
        - healthx
      initialDelaySeconds: 5
      periodSeconds: 5
```
#### grpc程序
```yaml
livenessProbe:
      exec:
        command:
        - /main
        - grpc
        - healthx
      initialDelaySeconds: 5
      periodSeconds: 5
```




