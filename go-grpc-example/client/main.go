package main

import (
	"context"
	"fmt"
	pb "github.com/grpcclient/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	//连接到server端，此处禁用安全传输，没有加密和验证
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	fmt.Println("connect success")
	defer conn.Close()

	client := pb.NewSayHelloClient(conn)
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "feng"})
	fmt.Println(resp.GetResponseMsg())
	fmt.Println("end...")
}
