package main

import (
	"context"
	"fmt"
	pb "github.com/grpcinterceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"runtime"
	"time"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("./key/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	conn, err := grpc.Dial("127.0.0.1:9097", grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(LogUnaryClientIntercrptor()))
	if err != nil {
		log.Fatalf("creds failure", err)
	}

	fmt.Println("connect success ")
	defer conn.Close()

	client := pb.NewHelloClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "feng"})
	if err != nil {
		fmt.Println("resp", err)
	}
	fmt.Println(resp.GetMessage())
}

func LogUnaryClientIntercrptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 预处理 (pre-processing)
		start := time.Now()

		//获取正在运行程序的操作系统
		cos := runtime.GOOS
		ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)

		err := invoker(ctx, method, req, reply, cc, opts...)

		//后处理
		end := time.Now()
		log.Printf("RPC:%s, client-OS: %v req: %v start time: %s , end time : %s err: %v", method, cos, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
		return err
	}

}
