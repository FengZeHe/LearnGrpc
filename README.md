# Learn gRPC



## protobuf

### 什么是protobuf
protocol buffer是一种序列化数据的方法，该数据可以通过有线传输或存储在文件中。JSON和XML等格式也用于序列化数据。protobuf是gRPC中序列化结构化数据的默认方法。和JSON、XML一样，Protobuf是与语言和平台无关的，protobuf仅专注于尽可能快的序列化和饭序列化数据的能力，另一个重要优化方法是通过传输数据尽可能小占用网络带宽。

### 什么是序列化工具
序列化：将结构数据或对象转换成能够被存储和传输（如网络传输）的格式，同时应当要保证这个序列化结果在之后（另外一个计算环境中）能够被重建回原来的结构数据或对象。
### .proto文件的用途
用于表示序列化数据的定义，包含称为消息的配置，可以编译原始文件以使用用户的编程语言生成代码。

### 安装protobuf
#### Mac下


#### ubuntu下


### 进行序列化与反序列化
```
    Please specify a program using absolute path or make sure the program is available in your PATH system variable
    --gofast_out: protoc-gen-gofast: Plugin failed with status code 1.
```

### protobuf 底层协议