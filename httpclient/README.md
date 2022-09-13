# http client 模块介绍
> 这是一个 HTTP 客户端的扩展和封装, 对外提供统一的http调用

## 预警
1. 现在http client和hystrix默认超时都是2s, 如果使用方业务比较耗时, 请设置较长的超时时间

## 描述
1. http client 是在http原生client的封装
2. http client 提供了强大的Middleware机制, 让扩展功能更加简单方便
3. http client提供了重试机制和超时设置
4. http client提供option的机制用于灵活配置

## 默认值
1. 默认超时时间是2s
2. 默认超时重试机制 `NewRetrier(NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))`
3. hystrix默认最大超时2s
4. hystrix默认最大并发请求5000
5. hystrix默认错误百分比阈值10
6. hystrix默认Sleep Window 10
7. hystrix默认名字 http.client

## 使用

```go
import "github.com/gowins/dionysus/httpclient"
```

### 简单的GET请求

```go
client := httpclient.New(httpclient.WithHTTPTimeout(100 * time.Millisecond))
rsp, err := client.Get(URL, nil)
if err != nil{
	panic(err)
}

// rsp是标准的*http.Response对象
body, _ := ioutil.ReadAll(res.Body)
fmt.Println(string(body))
```

### 简单的Do请求

```go
client := httpclient.New(httpclient.WithHTTPTimeout(100 * time.Millisecond))
req, _ := http.NewRequest(http.MethodGet, URL, nil)
rsp, err := client.Do(req)
if err != nil {
	panic(err)
}

// rsp是标准的*http.Response对象
body, _ := ioutil.ReadAll(res.Body)
fmt.Println(string(body))
```

### 超时重试设置

```go
client := httpclient.New(
    httpclient.WithHTTPTimeout(100 * time.Millisecond),
    httpclient.WithRetryCount(2),
    httpclient.WithRetrier(httpclient.NewRetrier(httpclient.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
)
```

## 使用hystrix中间件

```go
client := httpclient.New(
    httpclient.WithHTTPTimeout(100 * time.Millisecond),
    httpclient.WithRetryCount(2),
    httpclient.WithRetrier(httpclient.NewRetrier(httpclient.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
    httpclient.WithMiddleware(httphystrix.Middleware(
    		httphystrix.WithCommandName("MyCommand"),
    		httphystrix.WithMaxConcurrentRequests(100),
    		httphystrix.WithErrorPercentThreshold(25),
    		httphystrix.WithSleepWindow(10),
    		httphystrix.WithRequestVolumeThreshold(10),
    	)),
)
```

## transport设置
1. http client默认支持的是http.DefaultTransport
Client.Transports属性中包含：
    MaxIdleConns  所有host的连接池最大连接数量，默认无穷大
    MaxIdleConnsPerHost  每个host的连接池最大空闲连接数,默认2
    MaxConnsPerHost 对每个host的最大连接数量，0表示不限制
    如果MaxConnsPerHost=1，则只有一个http client被创建.
    如果MaxIdleConnsPerHost=1，则会缓存一个http client.

```go
// http client 默认的配置
&http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second, //限制建立TCP连接的时间
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second, //限制TLS握手的时间
    ExpectContinueTimeout: 1 * time.Second, //限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待。
    MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1, // http client缓存数量
}
```

2. 设置和替换自己的transport

```go
httpclient.New(httpclient.WithTransport("自定义transport"))
httpclient.Clone(httpclient.WithTransport("自定义transport"))

// 在使用的时候可以通过option的方式设置自己的transport
// 在调用Clone方法时, 直接复用父client的transport
```