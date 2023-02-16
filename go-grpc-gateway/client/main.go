package main

import (
	"context"
	"fmt"
	pb "github.com/grpcgateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9099", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed", err)
	}

	client := pb.NewHelloClient(conn)
	req, err := client.Sayhello(context.Background(), &pb.HelloRequest{RequestName: "Crazy Thursday!"})
	if err != nil {
		log.Fatalf("gRPC call sayhello failed", err)
	}
	fmt.Printf(req.ResponseMsg)

}
