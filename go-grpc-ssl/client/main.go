package main

import (
	"context"
	"fmt"
	pb "github.com/grpcssl/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("./conf/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println("Connection failure ", err)
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
