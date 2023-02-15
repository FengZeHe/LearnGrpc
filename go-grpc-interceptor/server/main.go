package main

import (
	"context"
	"fmt"
	pb "github.com/grpcinterceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"time"
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

	grpcserver := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LogUnaryServerInterceptor()))
	//grpcserver := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LoggingInterceptor), grpc.UnaryInterceptor(LogUnaryServerInterceptor()))
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
	fmt.Printf("gRPC method: %s, %v \n", info.FullMethod, resp)
	return resp, err
}

func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 预处理
		start := time.Now()
		//从传入的上下文获取数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("couldn't parse incoming context metadata")
		}
		//检索客户端操作系统 ，如果不存在则为空
		os := md.Get("client-os")
		// 获取客户端ip地址
		ip, err := getClientIP(ctx)
		if err != nil {
			return nil, err
		}

		//RPC方法真正的执行逻辑 ,调用RPC方法(invoking RPC method)
		m, err := handler(ctx, req)
		end := time.Now()
		log.Printf("RPC: %s,client-OS: '%v' and IP: '%v' req:%v start time: %s, end time: %s, err: %v", info.FullMethod, os, ip, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
		return m, err
	}
}

func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return " ", fmt.Errorf("can not parse client ip address")
	}
	return p.Addr.String(), nil

}
