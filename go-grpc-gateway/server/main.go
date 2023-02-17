package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/grpcgateway/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

type server struct {
	pb.UnimplementedHelloServer
}

const (
	Address string = ":9099"
	Network string = "tcp"
)

func main() {
	listen, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("listen failed", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &server{})

	go func() {
		fmt.Printf("gRPC server start..")
		err = grpcServer.Serve(listen)
		if err != nil {
			log.Fatalf("gRPC Server error", err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:9099",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	err = pb.RegisterHelloHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func (s *server) Sayhello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("hi %s FROM Server", req.RequestName)
	return &pb.HelloResponse{ResponseMsg: "hi " + req.RequestName}, nil
}
