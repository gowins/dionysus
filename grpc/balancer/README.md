## grpc client balancer
### 设计目标
       负载均衡的实现方式有很多，可以是单独的硬件设备也可以是单独的进程。grpc原生实现了loadbalance，这里我们借助这一设计完成客户端负载均衡能力。
    基于此可以不用维护单独的连接池，简化使用的难度并充分使用grpc的原生能力。
##### 参数说明
    endPoints 需要连接的服务地址，可以是多个相同的地址，增加负载能力。
#### 使用说明
   ```golang
 	conf := &clientbalancer.GrpcClientConfig{
 		EndPoints: []string{"127.0.0.1:8087"},
 	} 
 
 	con, err := clientbalancer.NewClient(conf)
 	if err != nil{
 		panic(err.Error())
 		return
 	}
 	c := pb.NewGreeterClient(con.Client)
 	for {
 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 		if resp, err:=c.SayHello(ctx, &pb.HelloRequest{Name: "test"}); err !=nil{
 			log.Error(err.Error())
 		}else {
 			fmt.Println(resp.Message)
 		}
 		cancel()
 		time.Sleep(200*time.Millisecond)
 	}
   ```

#### grpc balance 实现原理详解
![实现原理](../../../docs/images/grpcbalance.jpg)