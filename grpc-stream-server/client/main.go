package main

import (
	"context"
	pb "github.com/grpcserverstarem/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

const Address string = ":9093"

var grpcClient pb.StreamServerClient

func listValue() {
	req := pb.SimpleRequest{
		Data: "stream server grpc ",
	}
	stream, err := grpcClient.ListValue(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call listvalue error:", err)
	}

	for {
		res, err := stream.Recv()

		// 判断消息流是否已经结束
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream error :", err)
		}
		log.Println(res.StreamValue)
	}
}

func main() {
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net connect error:", err)
	}
	defer conn.Close()

	grpcClient = pb.NewStreamServerClient(conn)
	listValue()
}
