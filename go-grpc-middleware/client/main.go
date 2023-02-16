package main

import (
	"context"
	"fmt"
	pb "github.com/grpcmiddleware/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"runtime"
	"time"
)

var client pb.HelloClient

func main() {
	conn, err := grpc.Dial("127.0.0.1:9098", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(LogUnaryInterceptor()))
	if err != nil {
		log.Fatalf("error", err)
	}
	defer conn.Close()
	client = pb.NewHelloClient(conn)
	resp, err := client.Sayhello(context.Background(), &pb.HelloRequest{RequestName: "six six six"})
	if err != nil {
		log.Fatalf("rpc methods failed :", err)
	}
	fmt.Printf(resp.ResponseMsg)

}

func LogUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		//预处理
		start := time.Now()

		cos := runtime.GOOS
		ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)
		err := invoker(ctx, method, req, reply, cc, opts...)
		end := time.Now()

		log.Fatalf(" RPC %s clientOS :%v ,start time: %v,end time: %v", method, cos, start.Format(time.RFC3339), end.Format(time.RFC3339))
		return err

	}
}
