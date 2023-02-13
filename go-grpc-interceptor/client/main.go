package main

import (
	"context"
	"fmt"
	pb "github.com/grpcinterceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("./key/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	conn, err := grpc.Dial("127.0.0.1:9097", grpc.WithTransportCredentials(creds))
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
