# registry

1. registry的定义参考了micro registry, 同时默认实现了etcdv3

## 注册中心加载方式

### 环境变量

`GAPI_REGISTRY=etcdv3://127.0.0.1:2379&timeout=1s`
// timeout 连接超时时间


### flag

`./main grpc -r etcdv3://127.0.0.1:2379&timeout=1s`
-r 参数指定注册中心的地址

### 显示调用初始化
import "github.com/gowins/dionysus/grpc/registry"
`registry.Init("etcdv3://127.0.0.1:2379&timeout=1s")`

### 总结
1. 以上三种方式都可以对注册中心进行初始化


## 注册中心接口介绍

```go
type Registry interface {
    // 用于初始化调用
	Init(opts ...Option) error
	// 注册服务信息
	Register(*Service, ...RegisterOption) error
    // 服务信息反注册
	Deregister(*Service) error
	// 根据服务服务名字获取服务信息
	GetService(string) ([]*Service, error)
	// 列出所有的服务
	ListServices() ([]*Service, error)
	// 开启服务监控模式, 实时监控服务的变更
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Watcher interface {
	// Next是一个阻塞方法, 当服务有更新的时候才会进行下一步
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Action  string
	Service *Service
}
```
