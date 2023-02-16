package main

import (
	"context"
	"fmt"
	pb "github.com/grpcgateway/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedHelloServer
}

const (
	Address string = ":9099"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("listen failed", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &server{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("gRPC Server error", err)
	}
}

func (s *server) Sayhello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("hi %s FROM Server", req.RequestName)
	return &pb.HelloResponse{ResponseMsg: "hi " + req.RequestName}, nil
}
