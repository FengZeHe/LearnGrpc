package main

import (
	pb "github.com/grpcstreamclient/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	// Address 监听地址
	Address string = ":9094"
	// Network 网络通信协议
	Network string = "tcp"
)

type SimpleService struct {
	pb.UnimplementedStreamClientServer
}

func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	//从流中获取消息
	for {
		res, err := srv.Recv()
		if err == io.EOF {
			return srv.SendAndClose(&pb.SimpleResponse{Value: "ok"})
		}
		if err != nil {
			return err
		}

		log.Println(res.StreamData)
	}

}

func main() {
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("listen port faield", err)
	}

	//新建gRPC服务器实例
	grpcServer := grpc.NewServer()

	//在gRPC服务器注册服务
	pb.RegisterStreamClientServer(grpcServer, &SimpleService{})

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("grpcServer error ", err)
	}

}
