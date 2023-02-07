package main

import (
	pb "github.com/grpcserverstarem/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

const (
	Address string = ":9093"
	Network string = "tcp"
)

type StreamServer struct {
	pb.UnimplementedStreamServerServer
}

func (s *StreamServer) ListValue(req *pb.SimpleRequest, srv pb.StreamServer_ListValueServer) error {
	for n := 0; n < 5; n++ {
		err := srv.Send(&pb.StreamResponse{
			StreamValue: req.Data + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen error :", err)
	}
	log.Println(Address + " net.Listing...")
	grpcServer := grpc.NewServer()
	pb.RegisterStreamServerServer(grpcServer, &StreamServer{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer error", err)
	}
}
