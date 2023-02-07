package main

import (
	"context"
	pb "github.com/grpcstreamclient/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

const Address string = ":9094"

var streamClient pb.StreamClientClient

func routeList() {
	stream, err := streamClient.RouteList(context.Background())
	if err != nil {
		log.Fatalf("upload list err", err)
	}
	for n := 0; n < 5; n++ {
		err = stream.Send(&pb.StreamRequest{StreamData: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err:", err)
		}
	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("routelist get response err", err)
	}
	log.Println(res.Value)
}

func main() {
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	defer conn.Close()

	//	建立gRPC连接
	streamClient = pb.NewStreamClientClient(conn)
	routeList()
}
