# protobuf ＆ gRPC笔记

### protobuf 是什么、能解决什么问题？

​	Protobuf(protocol buffer)是一种序列化数据的方法，该数据可以通过有线传输或存储在文件中。JSON和XML等格式也用于序列化数据。protobuf是gRPC中序列化结构化数据的默认方法;和JSON、XML一样，Protobuf是和语言及平台无关的，protobuf仅专注于尽可能快的序列化和反序列化数据的能力，另一个重要优化方案是通过传输数据尽可能小占用网络带宽。

### 编写protobuf(.proto)文件

​	proto文件用于表示序列化数据的定义，包含消息的配置，可以编译原始文件让各种编程语言生成代码。

​	protobuf里有几个关键字：`syntax`表示protobuf版本，`package` 表示包名，option go_package表示文件输出的位置(路径)，这两条指令能合在一起写： `option go_package="./path;packageName";message关键字表示消息，既需要传输的数据格式定义，类似于go的struct。下面定义了几个字段，最后跟着数字的叫字段标识，在同一个message中不能重复，并且是[0，2^29-1]范围内的整数。service定义服务

```protobuf
syntax = "proto3";
package person;
option go_package="./person";
message Person{
  string person_name = 1;
  string person_address = 2;
  int64 person_age = 3;
  int64 person_id = 4;
  }
```



### 安装protobuf

一般来说有这几种安装方式：

#### Mac

```shell
brew install protobuf 
brew install protoc-gen-go
//验证
protoc --version
```



#### ubuntu

```shell
sudo apt update
sudo apt install protobuf-compiler
//验证
protoc --version
```



#### 编译安装

```shell
wget https://github.com/google/protobuf/releases/download/v3.5.1/protobuf-all-3.5.1.zip

unzip protobuf-all-3.5.1.zip

cd protobuf-3.5.1/

./configure

make

make install

//配置环境变量
cd ~
vim .bash_profile

export PROTOBUF=/usr/local/protobuf 
export PATH=$PROTOBUF/bin:$PATH
// 使配置生效
source .bash_profile
```



### go 实践

既然protobuf是用来序列化和反序列化的工具，那我们用protobuf生成对象序列化再反序列化。

看一眼项目结构是这样的：

```
.
├── go.mod
├── go.sum
├── main.go
├── person
│   └── person.pb.go
└── proto
    └── person.proto
```

切换到proto文件夹下，使用指令` protoc --go_out=..  *.proto`，在person文件夹生成文件`person.pb.go`

下载需要的包：

```
go get -u google.golang.org/protobuf
```

``` go
// 将一个对象序列化再反序列化
func main() {
	person1 := person.Person{}
	person1.PersonName = "feng"
	person1.PersonAge = 18
	person1.PersonAddress = "China"
	person1.PersonId = 1

	bytes, _ := json.Marshal(person1)
	fmt.Println("bytes =", bytes, " \n", "person1 : ", string(bytes))

	out, err := proto.Marshal(&person1)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("out = ", out)

	person2 := person.Person{}
	err = proto.Unmarshal(out, &person2)

	if err != nil {
		fmt.Println(err)
	}
	bytes, _ = json.Marshal(person2)
	fmt.Println("person2", string(bytes))
}
```

完整的代码在这里 `https://github.com/FengZeHe/LearngRPC/tree/main/go-protoc-example`



#### 可能会遇到的问题

1. 你的Mac没有装`homebrew`，但是下载链接又超级慢。这时候可以使用这一条国内专用指令安装homebrew。

   ```
   /bin/zsh -c "$(curl -fsSL https://gitee.com/cunkai/HomebrewCN/raw/master/Homebrew.sh)"
   ```

2. `protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC See https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code for more information.`

   `github.com/golang/protobuf/protoc-gen-go `这个包已经过时，要使用`google.golang.org/protobuf`这个包



## gRPC概述

### 什么是RPC

​	RPC（remote procedure call）远程过程调用协议，是一种通过网络从远程计算机上请求服务。大致就是用了以后调用远程服务就像调用本机服务一样方便（服务之间发生进程间调用）。

### 什么是gRPC

​	gRPC是一个现代开源高性能远程过程调用（RPC）框架，能在任何环境中运行。它能有效连接数据中心内和跨数据中心的服务，支持负载均衡、跟踪、监控检查和身份认证。将设备、移动端程序和浏览器连接到后端服务。



## gRPC实践

我们来实现一个客户端client调用服务端server里 sayHello的函数的例子

#### protobuf文件编写

```protobuf
// 定义版本
syntax = "proto3";
option go_package =  ".;service";

// 定义服务
service SayHello {
  rpc SayHello(HelloRequest) returns (HelloResponse);
}

// 定义消息
message HelloRequest{
  string requestName = 1;
}

message  HelloResponse {
  string responseMsg = 1;
}
```



#### 服务端编写

服务端主要实现以下几个步骤：

- 创建gRPC Server对象

  ``` go
  grpcServer := grpc.NewServer()
  ```

- 将server(其包含需要被调用到的服务端接口)注册到gRPC Serve的内部注册中心；当接受到请求时，通过内部的服务发现，发现该服务端接口并进行逻辑处理

  ```
  pb.RegisterSayHelloServer(grpcServer, &server{})
  ```

- 创建Listen 监听TCP端口

  ```
  listen, _ := net.Listen("tcp", ":9091")
  ```

- gRPC Server 运行

  ``` go
  grpcServer.Serve(listen)
  ```

完整的服务端代码：

#### 客户端编写

- 创建与改定目标（服务端）的连接交互

  ```
  conn, err := grpc.Dial("127.0.0.1:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
  //此处禁用了安全传输，没有加密和验证
  ```

- 创建server的客户端对象

  ``` go
  client := pb.NewSayHelloClient(conn)
  ```

- 发送RPC请求，等待同步响应，得到回调后返回响应结果

  ``` go
  resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "feng"})
  	if err != nil {
  		fmt.Println("resp", err)
  }
  ```

- 输出响应结果

  ```
  fmt.Println(resp.GetResponseMsg())
  ```

完整代码在这里： `https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-example`





## gRPC的安全传输

#### TLS  概述＆握手过程

​	SSL和TLS协议可以为通信双方提供识别和认证的通道，从而确保通信的保密性和数据的完整性。在TLS握手过程中，通信双方交换消息以验证通信，互相确认并建立它们所要使用的加密算法以及会话密钥。

​	TLS在握手过程中能确定以下事情：

1. 确定双方通信所使用的TLS版本。

2. 确定双方所需要使用的密码组合。

3. 客户端通过服务器的公钥和数字证书上的签名验证服务端身份

4. 生成会话密钥，该密钥将用于握手结束后的对称加密。

   TLS握手详细过程是这样的：

5. "client hello"消息：客户端通过发送"client hello"消息向服务器发起握手请求，该消息包含了客户端所支持的TLS版本和密码组合供服务器选择，还有一个"client random"随机的字符串

6. "server hello"消息：服务器发送"server hello" 消息对客户端进行回应，该消息包含了数字证书，服务器选择的密码组合和"server random"随机字符串

7. 验证：客户端对服务器发来的证书进行验证，确保对方的合法身份

8. "Premaster secret" 字符串：客户端向服务器发送另一个随机字符串"premaster secret(预主密钥)"，这个字符串是经过服务器公钥加密的，因此只有对应的私钥能解密。

9. 生成私钥：服务器使用私钥解密"premaster secret"

10. 生成共享密钥：客户端和服务端均使用client random ,server random和premaster scret，并通过相同的算法生成共享密钥KEY

11. 客户端就绪：客户端发送经过共享密钥KEY加密过的"finished"信号

12. 服务器就绪：服务器发送经过共享密钥KEY加密过的"finished"信号

13. 达成安全通信：握手完成，上方使用对称加密进行安全通信。

#### TLS证书认证

生成证书的配置文件ca.conf和server.conf

ca.conf

```
[ req ]
default_bits       = 4096
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = GB
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GuangDong
localityName                = Locality Name (eg, city)
localityName_default        = Foshan
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = Foshan
commonName_max              = 64
commonName_default          = Foshan
```

server.conf

```
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GuangDong
localityName                = Locality Name (eg, city)
localityName_default        = Foshan
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = CommonName (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = Foshan

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1   = go-grpc-example #这里要指定好 
IP      = 127.0.0.1
```



#### 生成证书

切换到conf目录下

##### 生成CA根证书

1. 生成ca私钥，得到ca.key

   ```
   openssl genrsa -out ca.key 4096
   ```

2. 生成ca证书签发请求，得到ca.csr

   ```
   openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
   ```

   openssl req：生成自签名证书，-new 指生成证书请求、-sha256 指使用 sha256 加密、-key 指定私钥文件、-x509 指输出证书、-days 3650 为有效期，-config 指定配置文件

3. 生成ca根证书，得到ca.crt

   ```
   openssl x509 -req -days 3650 -in ca.csr -signkey ca.key -out ca.crt
   ```



##### 生成终端用户证书

1. 生成私钥，得到server.key

   ```
   openssl genrsa -out server.key 4096
   ```

2. 生成证书签发请求，得到server.csr

   ```
   openssl req -new -sha256 -out server.csr -key server.key -config server.conf
   ```

3. 用CA证书生成终端用户证书，得到server.crt

   ```
   openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
   ```



##### Server端

1. 根据服务端引用的证书文件和密钥构造TLS凭证。

   ```
   // 跟禁用加密的server端区别在这里
   creds, err := credentials.NewServerTLSFromFile("./conf/server.pem", "./conf/server.key")
   
   grpcServer := grpc.NewServer(grpc.Creds(creds))
   ```



##### Client 端

1. 客户端引用证书文件和密钥构造TLS凭证

   ``` go
   creds, err := credentials.NewClientTLSFromFile("./conf/server.pem", "go-grpc-example")
   ```

2. grpc.Dial 配置连接选项

   ``` go
   conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(creds))
   ```

#### 目录结构

```
.
├── client
│   ├── main.go
│   └── proto
│       ├── hello.pb.go
│       ├── hello.proto
│       └── hello_grpc.pb.go
├── conf
│   ├── ca.conf
│   └── server.conf
├── go.mod
├── go.sum
└── server
    ├── main.go
    └── proto
        ├── hello.pb.go
        ├── hello.proto
        └── hello_grpc.pb.go
```

完整代码在这里：https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-ssl



#### 自定义Token认证

##### 实现步骤

1. 客户端请求时带上Credentials
2. 服务端取出Credentials并验证有效性(一般配合拦截器使用)。

##### 编写proto文件

```protobuf
syntax = "proto3";

option go_package=".;grpctoken";

message SimpleRequest {
  string data = 1;
}

message SimpleResponse {
  int32 code = 1;
  string value =2;
}

service SayHello{
  rpc SayHello(SimpleRequest) returns (SimpleResponse){};
}
```

##### 生成Go代码

```
protoc --go_out =. ./*proto
protoc --go-grpc_out=. /*proto
```

#### 生成自签证书

##### 生成CA根证书

1. 生成ca私钥，得到ca.key

   ```
   openssl genrsa -out ca.key 4096
   ```

2. 生成ca证书签发请求，得到ca.csr

   ```
   openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
   ```

   openssl req：生成自签名证书，-new 指生成证书请求、-sha256 指使用 sha256 加密、-key 指定私钥文件、-x509 指输出证书、-days 3650 为有效期，-config 指定配置文件

3. 生成ca根证书，得到ca.crt

   ```
   openssl x509 -req -days 3650 -in ca.csr -signkey ca.key -out ca.crt
   ```

##### 生成终端用户证书

1. 生成私钥，得到server.key

   ```
   openssl genrsa -out server.key 4096
   ```

2. 生成证书签发请求，得到server.csr

   ```
   openssl req -new -sha256 -out server.csr -key server.key -config server.conf
   ```

3. 用CA证书生成终端用户证书，得到server.crt

   ```
   openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
   ```



#### Auth

``` go
// Token 认证
type Token struct {
	AppID     string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppID, "app_secret": t.AppSecret}, nil
}

// RequireTarnsportSecurity 是否要基于TLS认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}
```



#### Server端

``` go
func main() {
	//	监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("lieten error", err)
	}
	// 从引用证书文件和密钥文件为服务构造TLS凭证
	creds, err := credentials.NewServerTLSFromFile("./pkg/tls/server.pem", "./pkg/tls/server.key")
	if err != nil {
		log.Fatalf("Failed to grnerate credentials %v", err)
	}

	//一元拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = Check(ctx)
		if err != nil {
			return
		}
		return handler(ctx, req)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptor))
	pb.RegisterSayHelloServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpc server error %v", err)
	}
}


func (s *server) SayHello(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	res := pb.SimpleResponse{
		Code:  200,
		Value: "hello" + req.Data,
	}
	return &res, nil
}


// check 验证token
func Check(ctx context.Context) error {
	//从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "获取token失败")
	}
	var (
		appID     string
		appSecret string
	)
	if value, ok := md["app_id"]; ok {
		appID = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appID != "grpc_token" || appSecret != "12345678" {
		return status.Errorf(codes.Unauthenticated, "Token无效 %v %v ", appID, appSecret)
	}
	return nil
}
```

#### Client端

```
//引用证书文件
creds, err := credentials.NewClientTLSFromFile("./pkg/tls/server.pem", "go-grpc-example")

token := auth.Token{AppID: "grpc_token", AppSecret: "12345678"}

//连接到服务器
conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
```

##### 目录结构

```
.
├── client
│   └── main.go
├── go.mod
├── pkg
│   ├── auth
│   │   └── auth.go
│   └── tls
│       ├── server.key
│       └── server.pem
├── proto
│   └── hello.proto
└── server
    └── main.go
```

完成代码在这里：`https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-token`

### 可能遇到的问题

`code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs instead"`

​	出现这种问题的原因是使用的证书没有开启SAN扩展。SAN（Subject Alternative Name）是SSL标准X509中定义的一个扩展，使用该字段生成的SSL证书扩展支持的域名，能使一个证书支持多个不同域名解析。



### gRPC 拦截器

#### 概述

gRPC提供了拦截器(Interceptor)功能，包括客户端拦截器和服务端拦截器。可以在接收到请求或者发起请求之前优先对氢气中的数据做一些处理后再转交给指定服务处理相应；很适合做处理验证、日志等流程。

#### 拦截器有哪些类型

- `UnaryServerInterceptor` 服务端拦截，在服务端接收请求的时候进行拦截。

- `UnaryClientInterceptor` 客户端拦截器，在客户端真正发起调用之前，进行拦截。

- `StreamClientInterceptor` 在流式客户端调用时，通过拦截 clientstream 的创建，返回一个自定义的 clientstream, 可以做一些额外的操作。

- `StreamServerInterceptor` 在服务端接收到流式请求的时候进行拦截。



##### 客户端一元拦截器(Client Interceptor)

​	客户端的一元拦截器类型为`UnaryClientInterceptor`,实现分为`预处理(pre-poressing)`、`调用RPC方法(invoking RPC method)`和`后处理(post-processing)`三个阶段。

​	参数含义如下：

- `ctx` ：Go语言中的上下文，一般和Goroutine配合使用，起到超时控制的效果。
- `method`: 当前调用的RPC方法名
- `req` ： 本次请求的参数，只有在**处理前**阶段修改才有效
- `reply` ： 本次请求响应，需要在**处理后** 阶段才能获得
- `cc` : gRPC连接信息
- `invoker` ： 可以看做是当前RPC方法，一般在拦截器中调用invoker能达到调用RPC方法的效果，底层也是RPC处理
- `opts` ：本次调用指定的options信息

``` go
type UnaryClientInterceptor func(
    ctx context.Context, 
    method string, 
    req, 
    reply interface{}, 
    cc *ClientConn, 
    invoker UnaryInvoker, 
    opts ...CallOption,
) error
```



##### 客户端流拦截器 (Stream Interceptor)

​	客户端流拦截器的实现包括预处理和流操作拦截，并不能在事后进行RPC方法调用和后处理，而是拦截用户对流的操作。

``` go
type StreamClientInterceptor func(
    ctx context.Context, 
    desc *StreamDesc, 
    cc *ClientConn, 
    method string, 
    streamer Streamer, 
    opts ...CallOption,
) (ClientStream, error)
```



##### 服务端一元拦截器

​	服务端一元拦截器类型为`UnaryServerInterceptor` ，一共包含4个参数，包括RPC上下文、RPC请求参数、RPC方法的所有信息、RPC方法本身。

``` go
type UnaryServerInterceptor func(
    ctx context.Context, 
    req interface{}, 
    info *UnaryServerInfo, 
    handler UnaryHandler,
) (resp interface{}, err error)
```



##### 服务端流拦截器

​	服务端流拦截器类型为`StreamServerInterceptor`，





## gRPC流

### 概述

​	当数据量大或者需要不断传输数据的时候，就应该使用流式RPC，它允许我们一边处理一边传输数据。流式RPC分为服务端流式RPC和客户端流式RPC。**服务端流式RPC过程是：客户端发送请求到服务器，拿到一个流读取返回的消息队列。客户端读取返回的流，直到里面没有任何消息。客户端流式RPC的过程是：客户端不断向服务端发送数据流，在发送结束后由服务端返回一个响应。**

### 实践

#### 实现服务端流RPC

##### 定义proto文件

``` protobuf
syntax ="proto3";
option go_package = ".;StreamServer";
// 定义发送请求消息
message SimpleRequest{
  string data = 1;
}
// 定义流式相应消息
message StreamResponse{
  string stream_value = 1;
}

// 定义服务方法ListValue
service StreamServer {
	// 流式服务端RPC，因此在returns的参数天 stream
  rpc ListValue(SimpleRequest) returns(stream StreamResponse){};
}
```

##### 编译proto文件

```
protoc --go_out=. *.proto
protoc --go-grpc_out=. *.proto
```

##### 编写Server端的程序

- 主要是实现定义的ListValue方法

``` go
func (s *StreamServer) ListValue(req *pb.SimpleRequest, srv pb.StreamServer_ListValueServer) error {
	for n := 0; n < 5; n++ {
		err := srv.Send(&pb.StreamResponse{
			StreamValue: req.Data + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
```

##### 编写Client端的程序

``` go
func listValue() {
	req := pb.SimpleRequest{
		Data: "stream server grpc ",
	}
	stream, err := grpcClient.ListValue(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call listvalue error:", err)
	}

	for {
		res, err := stream.Recv()

		// 判断消息流是否已经结束
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream error :", err)
		}
		log.Println(res.StreamValue)
	}
}
```

完整代码在这里 `https://github.com/FengZeHe/LearngRPC/tree/main/grpc-stream-server`

#### 实现客户端流RPC

##### 编写proto文件

``` protobuf
syntax = "proto3";
option  go_package = ".;ClientStream";


message StreamRequest{
  string stream_data = 1;
}

message SimpleResponse {
  int32 code = 1;
  string value = 2;
}

service StreamClient{
  rpc RouteList(stream StreamRequest) returns ( SimpleResponse){};
}
```

##### 编写服务端程序

``` go
func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	//从流中获取消息
	for {
		res, err := srv.Recv()
    //如果已经读完了
		if err == io.EOF {
			return srv.SendAndClose(&pb.SimpleResponse{Value: "ok"})
		}
		if err != nil {
			return err
		}

		log.Println(res.StreamData)
	}
}


func main(){
  ...
  //在grpc中注册服务
  pb.RegisterStreamClientServer(grpcServer, &SimpleService{})
  ...
}
```

##### 编写客户端程序

``` go
func routeList() {
	stream, err := streamClient.RouteList(context.Background())
	if err != nil {
		log.Fatalf("upload list err", err)
	}
	for n := 0; n < 5; n++ {
		err = stream.Send(&pb.StreamRequest{StreamData: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err:", err)
		}
	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("routelist get response err", err)
	}
	log.Println(res.Value)
}
```

​		完整代码在这里 `https://github.com/FengZeHe/LearngRPC/tree/main/grpc-stream-client`



#### 实现双向流式RPC

##### 概述

双向流式RPC的意思是客户端和服务端双方读写流发送消息序列，两个流单独操作，双方可以同时发送和接受消息。

##### 编写proto文件

``` protobuf
syntax = "proto3";
option go_package = ".;StreamConversations";

message StreamRequest {
     string question = 1;
}

message StreamResponse {
  string answer = 1;
}

service StreamConversations {
  rpc Conversations(stream StreamRequest) returns (stream StreamResponse){};
}
```

##### 编译文件

```
protoc --go_out=. *.proto
protoc --go-grpc_out=. *.proto
```

##### 编写服务端程序

``` go
//实现Conversations()方法
func (s *StreamService) Conversations(srv pb.StreamConversations_ConversationsServer) error {
	n := 1
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = srv.Send(&pb.StreamResponse{
			Answer: "from stream server answer: the " + strconv.Itoa(n) + " question is " + req.Question,
		})
		if err != nil {
			return err
		}
		n++
		log.Printf("from stream client question: %s", req.Question)
	}
}

func main(){
  ...
  //在gRPC中注册服务
	pb.RegisterStreamConversationsServer(grpcServer, &StreamService{})
  ...
}
```



##### 编写客户端程序

``` go
// 实现conversations方法
func conversations() {
	stream, err := streamClient.Conversations(context.Background())
	if err != nil {
		log.Fatalf("stream failure")
	}
	for n := 0; n < 5; n++ {
		err := stream.Send(&pb.StreamRequest{Question: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.Answer)
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("Conversations close stream err: %v", err)
	}
}


func main(){
  ...
  //建立gRPC连接
  streamClient = pb.NewStreamConversationsClient(conn)
  conversations()
  ...
}
```

完整代码在这里`https://github.com/FengZeHe/LearngRPC/tree/main/grpc-stream-conversations`




### gRPC 拦截器

#### 概述

gRPC提供了拦截器(Interceptor)功能，包括客户端拦截器和服务端拦截器。可以在接收到请求或者发起请求之前优先对氢气中的数据做一些处理后再转交给指定服务处理相应；很适合做处理验证、日志等流程。

#### 拦截器有哪些类型

- `UnaryServerInterceptor` 服务端拦截，在服务端接收请求的时候进行拦截。

- `UnaryClientInterceptor` 客户端拦截器，在客户端真正发起调用之前，进行拦截。

- `StreamClientInterceptor` 在流式客户端调用时，通过拦截 clientstream 的创建，返回一个自定义的 clientstream, 可以做一些额外的操作。

- `StreamServerInterceptor` 在服务端接收到流式请求的时候进行拦截。

### 客户端

##### 客户端一元拦截器(Client Interceptor)

​	客户端的一元拦截器类型为`UnaryClientInterceptor`,实现分为`预处理(pre-poressing)`、`调用RPC方法(invoking RPC method)`和`后处理(post-processing)`三个阶段。

​	参数含义如下：

- `ctx` ：Go语言中的上下文，一般和Goroutine配合使用，起到超时控制的效果。
- `method`: 当前调用的RPC方法名
- `req` ： 本次请求的参数，只有在**处理前**阶段修改才有效
- `reply` ： 本次请求响应，需要在**处理后** 阶段才能获得
- `cc` : gRPC连接信息
- `invoker` ： 可以看做是当前RPC方法，一般在拦截器中调用invoker能达到调用RPC方法的效果，底层也是RPC处理
- `opts` ：本次调用指定的options信息

``` go
type UnaryClientInterceptor func(
    ctx context.Context, 
    method string, 
    req, 
    reply interface{}, 
    cc *ClientConn, 
    invoker UnaryInvoker, 
    opts ...CallOption,
) error
```

##### 客户端流拦截器 (Stream Interceptor)

​	客户端流拦截器的实现包括预处理和流操作拦截，并不能在事后进行RPC方法调用和后处理，而是拦截用户对流的操作。拦截器的区别也体现在请求参数上，如req参数变成了streamer。拦截过程

``` go
type StreamClientInterceptor func(
    ctx context.Context, 
    desc *StreamDesc, 
    cc *ClientConn, 
    method string, 
    streamer Streamer, 
    opts ...CallOption,
) (ClientStream, error)
```

##### 异同

流式拦截器同样分为三个阶段：**预处理、调用RPC方法、后处理**。预处理阶段和一元拦截器类似，但后面两个阶段则不同；StreamAPI的请求和响应都是通过Stream进行传递的，更进一步是通过Streamer调用SendMsg和RecvMsg这两个方法获取的。然后Streamer又是低啊用RPC方法来获得，所以在流拦截器中我们可以对streamer进行包装，进而实现SendMsg和RecvMsg这两个方法。

### 服务端

##### 服务端一元拦截器

​	服务端一元拦截器类型为`UnaryServerInterceptor` ，一共包含4个参数，包括RPC上下文、RPC请求参数、RPC方法的所有信息、RPC方法真正执行的逻辑。

``` go
type UnaryServerInterceptor func(
    ctx context.Context, 
    req interface{}, 
    info *UnaryServerInfo, 
    handler UnaryHandler,
) (resp interface{}, err error)
```

##### 服务端流拦截器

​	服务端流拦截器类型为`StreamServerInterceptor`，

```
type StreamClientInterceptor func(
    ctx context.Context, 
    desc *StreamDesc, 
    cc *ClientConn, 
    method string, 
    streamer Streamer, 
    opts ...CallOption,
) (ClientStream, error)
```


#### 实践

##### 实现客户端和服务端的一元拦截器
- 客户端实现一元拦截器在`grpc.Dail()`指定即可
  ``` go
  //  实现客户端一元拦截器
    conn, err := grpc.Dial("127.0.0.1:9097", grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(LogUnaryClientIntercrptor()))
    
  // 定义客户端一元拦截器 LogUnaryClientInterceptor()
    func LogUnaryClientIntercrptor() grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        ...xxx
    }
  }
  ```
- 服务端实现一元拦截器在`grpc.NewServer()`指定即可
  ``` go
  //服务端实现一元拦截器
	 grpcserver := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LogUnaryServerInterceptor()))
    
  // 服务端定义一元拦截器 定义名字LogUnaryServerInterceptor的一元拦截器
  func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	 return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
     xxxx...
  }
  ```
- 完整的代码在这里 `https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-interceptor`
  


#### 实现客户端和服务端流拦截器

- 概述
流拦截器过程分为3个阶段：
1. 预处理(pre-processing)
2. 调用RPC方法(invoking RPC method)
3. 后处理(post-processing)

-  StreamAPI的请求和响应都是通过Stream进行传递的，进一步是通过Streamer调用SendMsg和RecvMsg这两个方法
-  Streamer是调用RPC方法获取的，所以在流拦截器中我们对Stream进行包装，实现SendMsg和RecvMsg这两个方法。

##### 客户端的流拦截器

- 同样需要在`grpc.Dial`中指定流拦截器

``` go
// wrappedStream  用于包装 grpc.ClientStream 结构体并拦截其对应的方法。
type wrappedStream struct {
   grpc.ClientStream
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// 实现RecvMsg方法
func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("Receive a message (Type: %T) at %v \n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// 实现SendMsg方法
func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("Send a message (Type: %T) at %v \n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

// streamInterceptor 一个简单的 stream interceptor 示例。
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	// 返回的是自定义的封装过的 stream
	return newWrappedStream(s), nil
}

```



##### 服务端的流拦截器

- 服务端实现流拦截器也是类似的

``` go
type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RevcMsg(m interface{}) (err error) {
	fmt.Printf("Receive a message (Type: %T) at %s ", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) (err error) {
	fmt.Printf("Send a message (Type %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	//包装 grpc.ServerStream 以替换 RecvMsg SendMsg这两个方法。
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v", err)
	}
	return err
}
```

完整代码在这里 `https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-stream-interceptor`

#### **拦截器执行过程**

**一元拦截器**

- 1）预处理
- 2）调用RPC方法
- 3）后处理

**流拦截器**

- 1）预处理
- 2）调用RPC方法 获取 Streamer
- 3）后处理
    - 调用 SendMsg 、RecvMsg 之前
    - 调用 SendMsg 、RecvMsg
    - 调用 SendMsg 、RecvMsg 之后



## Go gRPC Middleware

#### 概述

Go gRPC Middleware是一个gRPC中间件，提供了拦截器的链式功能，常用来做身份认证、日志记录、监控、客户端重连等等。

##### 项目地址

`https://github.com/grpc-ecosystem/go-grpc-middleware`

#### 实践

- 使用`Middleware`实现多个拦截器

  ``` go
  // 在项目中引用这个包
  import "github.com/grpc-ecosystem/go-grpc-middleware"
  
  // 在server端定义两个Log一元拦截器
  func LogUnaryIntercptor() grpc.UnaryServerInterceptor {
  	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
  		...
    }
  }
  
  func LogUnaryIntercptorTwo() grpc.UnaryServerInterceptor {
  	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
  		...
  	}
  }
  
  // 在main函数中加入这两个拦截器
  grpcServer := grpc.NewServer(
  		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
  			LogUnaryIntercptorTwo(),
  			LogUnaryIntercptor(),
  		)),
  	)
  ```



## go-grpc-gateway


#### 概述

gRPC-Gateway是一个protoc插件。它读取gRPC服务定义并生成一个反向代理服务器，该服务器将RESTFUL JSON API转换为gRPC；gRPC-Gateway 能同时提供 gRPC 和 RESTful 风格的 API。



#### 步骤

1. 写一个grpc服务器(编写proto文件+编译+客户端程序+服务端程序)

2. 引用两个文件 `annotations.proto` 和 `http.proto`

   从官方存储库googleapis下载： `https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto`

   `https://github.com/googleapis/googleapis/blob/master/google/api/http.proto`

   ```
   ├── google
   │   └── api
   │       ├── annotations.proto
   │       └── http.proto
   └── hello
       ├── hello.pb.go
       ├── hello.proto
       └── hello_grpc.pb.go
   ```

    - 在proto文件中导入 `anntotations.proto`

      ``` protobuf
      import "google/api/annotations.proto";
      ```

3. proto文件添加gRPC注释（注释定义gRPC服务映射到JSON请求和响应，使用protobuf 时每个RPC服务必须使用google.api.HTTP来注释定义HTTP定义和路径）

   ``` protobuf
   service Hello {
     rpc sayhello(HelloRequest) returns (HelloResponse){
     //添加注释
       option (google.api.http) = {
         get: "/v1/hello/sayhello"
       };
     };
   }
   ```

4. 再次编译

   ```
   protoc -I ./proto \
      --go_out ./proto --go_opt paths=source_relative \
      --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
      ./proto/hello/hello.proto
   ```

   编译过后会生成`hello.pb.gw.go`文件

   ```
   .
   ├── google
   │   └── api
   │       ├── annotations.proto
   │       └── http.proto
   └── hello
       ├── hello.pb.go
       ├── hello.pb.gw.go
       ├── hello.proto
       └── hello_grpc.pb.go
   ```

5. 在main函数中添加并提供gRPC-Gateway mux

   ``` go
   	func main(){
       xxxx...
       conn, err := grpc.DialContext(
         context.Background(),
         "0.0.0.0:9099",
         grpc.WithBlock(),
         grpc.WithTransportCredentials(insecure.NewCredentials()),
       )
       if err != nil {
         log.Fatalf("Failed to dial server:", err)
       }
       gwmux := runtime.NewServeMux()
       err = pb.RegisterHelloHandler(context.Background(), gwmux, conn)
       if err != nil {
         log.Fatalln("Failed to register gateway:", err)
       }
       
       xxx...
   	}
   ```

    - 用postman测试接口
   完整代码在这里：`https://github.com/FengZeHe/LearngRPC/tree/main/go-grpc-gateway`




