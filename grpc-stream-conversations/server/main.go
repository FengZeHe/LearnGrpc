package main

import (
	pb "github.com/grpcstreamconversations/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	Address string = ":9095"
	NetWork string = "tcp"
)

type StreamService struct {
	pb.UnimplementedStreamConversationsServer
}

func (s *StreamService) Conversations(srv pb.StreamConversations_ConversationsServer) error {
	n := 1
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = srv.Send(&pb.StreamResponse{
			Answer: "from stream server answer: the " + strconv.Itoa(n) + " question is " + req.Question,
		})
		if err != nil {
			return err
		}
		n++
		log.Printf("from stream client question: %s", req.Question)
	}
}

func main() {
	listener, err := net.Listen(NetWork, Address)
	if err != nil {
		log.Fatalf("listen failure", err)
	}
	//	启动gRPC服务
	grpcServer := grpc.NewServer()
	//在gRPC中注册服务
	pb.RegisterStreamConversationsServer(grpcServer, &StreamService{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer failure ！Error =>", err)
	}
}
