package main

import (
	"context"
	pb "github.com/grpctoken/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

const (
	Address string = ":9096"
	Network string = "tcp"
)

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
