package main

import (
	"context"
	"fmt"
	pb "github.com/grpcinterceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type service struct {
	pb.UnimplementedHelloServer
}

const (
	Network = ":9097"
	Address = "tcp"
)

func main() {
	listener, err := net.Listen(Address, Network)
	if err != nil {
		log.Fatalf("listen network error", err)
	}

	creds, err := credentials.NewServerTLSFromFile("./key/server.pem", "./key/server.key")

	grpcserver := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LoggingInterceptor))
	pb.RegisterHelloServer(grpcserver, &service{})

	grpcserver.Serve(listener)
}

func (s *service) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("hello %v", in.Name)
	return resp, nil
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("gRPC method: %s, %v \n", info.FullMethod, req)
	resp, err := handler(ctx, req)
	fmt.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}
