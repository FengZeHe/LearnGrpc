package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/grpcmiddleware/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedHelloServer
}

const (
	Address string = ":9098"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("listen error :", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			LogUnaryIntercptorTwo(),
			LogUnaryIntercptor(),
		)),
	)
	pb.RegisterHelloServer(grpcServer, &server{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("grpc server failed", err)
	}
}

func (s *server) Sayhello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{ResponseMsg: "Hello " + req.RequestName}, nil
}

func LogUnaryIntercptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 预处理
		start := time.Now()
		//传入上下文获取数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("can not parse incoming context metadata")
		}
		os := md.Get("client-os")

		//rpc方法执行逻辑，调用rpc方法
		m, err := handler(ctx, req)
		end := time.Now()
		log.Printf("Interceptor 1:RPC: %s ,client os : %v ,start time %s,end time %s", info.FullMethod, os, start.Format(time.RFC3339), end.Format(time.RFC3339))
		return m, err
	}
}

func LogUnaryIntercptorTwo() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 预处理
		start := time.Now()
		//传入上下文获取数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("can not parse incoming context metadata")
		}
		os := md.Get("client-os")

		//rpc方法执行逻辑，调用rpc方法
		m, err := handler(ctx, req)
		end := time.Now()
		log.Printf("Interceptor 2:RPC: %s ,client os : %v ,start time %s,end time %s", info.FullMethod, os, start.Format(time.RFC3339), end.Format(time.RFC3339))
		return m, err
	}
}
