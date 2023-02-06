package main

import (
	"context"
	"fmt"
	pb "github.com/grpcserver/proto"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil
}

func main() {
	//开启端口
	listen, _ := net.Listen("tcp", ":9091")
	//创建gRPC服务
	grpcServer := grpc.NewServer()
	//在grpc服务端中注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})

	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("failed to server", err)
		return
	}
}
