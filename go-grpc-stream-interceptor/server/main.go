package main

import (
	"fmt"
	pb "github.com/grpcstreaminterceptor/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type streaminterceptor struct {
	pb.UnimplementedStreamConversationsServer
}

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RevcMsg(m interface{}) (err error) {
	fmt.Printf("Receive a message (Type: %T) at %s ", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) (err error) {
	fmt.Printf("Send a message (Type %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	//包装 grpc.ServerStream 以替换 RecvMsg SendMsg这两个方法。
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v", err)
	}
	return err
}

const (
	Address string = ":9097"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("listen failure", err)
	}

	grpcServer := grpc.NewServer(grpc.StreamInterceptor(streamInterceptor))
	pb.RegisterStreamConversationsServer(grpcServer, &streaminterceptor{})

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Serve error", err)
	}
}

func (s *streaminterceptor) Converstaion(srv pb.StreamConversations_ConverstaionServer) error {
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
			Answer: "from stream server answer : the " + strconv.Itoa(n) + "question is " + req.Question,
		})
		n++
		log.Printf("from stream client question : %s", req.Question)
	}
}
