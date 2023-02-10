package main

import (
	"context"
	"github.com/grpctoken/pkg/auth"
	pb "github.com/grpctoken/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const Address string = ":9096"

var grpcClient pb.SayHelloClient

func main() {
	// 在引用的证书文件中 为客户端构造TLS凭证
	creds, err := credentials.NewClientTLSFromFile("./pkg/tls/server.pem", "go-grpc-example")

	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	token := auth.Token{AppID: "grpc_token", AppSecret: "12345678"}

	//	连接到服务器
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		log.Fatalf("net connect failure %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewSayHelloClient(conn)
	sayHello()
}

func sayHello() {
	req := pb.SimpleRequest{Data: " dawei "}

	res, err := grpcClient.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err %v", err)
	}
	log.Println(res)
}
