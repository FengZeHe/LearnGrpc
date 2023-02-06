package main

import (
	"context"
	"fmt"
	pb "github.com/grpcssl/server/proto"
	"google.golang.org/grpc"
	"net"
)

type service struct {
	pb.UnimplementedSayHelloServer
}

// SayHello
func (s *service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{ResponseName: "ssl - hello " + req.RequestName}, nil
}

func main() {
	//开启端口
	listenport, _ := net.Listen("tcp", ":9092")

	// 创建gRPC服务
	grpcServer := grpc.NewServer()

	//在注册中心注册服务
	pb.RegisterSayHelloServer(grpcServer, &service{})

	//	启动服务
	err := grpcServer.Serve(listenport)
	if err != nil {
		fmt.Println("error ", err)
		return
	}

}
