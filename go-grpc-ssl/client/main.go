package main

import (
	"context"
	"fmt"
	pb "github.com/grpcssl/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Connection failure ")
		return
	}

	defer conn.Close()

	//建立连接
	client := pb.NewSayHelloClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "feng"})
	if err != nil {
		fmt.Println("err !", err)
	}

	fmt.Println(resp.GetResponseName())
}
