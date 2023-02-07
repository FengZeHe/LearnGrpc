package main

import (
	"context"
	pb "github.com/grpcstreamconversations/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strconv"
)

const Address string = ":9095"

var streamClient pb.StreamConversationsClient

func conversations() {
	stream, err := streamClient.Conversations(context.Background())
	if err != nil {
		log.Fatalf("stream failure")
	}
	for n := 0; n < 5; n++ {
		err := stream.Send(&pb.StreamRequest{Question: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.Answer)
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("Conversations close stream err: %v", err)
	}
}

func main() {
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	defer conn.Close()

	//	建立gRPC连接
	streamClient = pb.NewStreamConversationsClient(conn)
	conversations()

}
