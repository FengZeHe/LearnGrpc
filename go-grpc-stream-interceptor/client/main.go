package main

import (
	"context"
	"fmt"
	pb "github.com/grpcstreaminterceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strconv"
	"time"
)

const Address string = "127.0.0.1:9097"

var client pb.StreamConversationsClient

// wrappedStream  用于包装 grpc.ClientStream 结构体并拦截其对应的方法。
type wrappedStream struct {
	grpc.ClientStream
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// 实现RecvMsg方法
func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("Receive a message (Type: %T) at %v \n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// 实现SendMsg方法
func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("Send a message (Type: %T) at %v \n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

// streamInterceptor 一个简单的 stream interceptor 示例。
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	// 返回的是自定义的封装过的 stream
	return newWrappedStream(s), nil
}

func main() {
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(streamInterceptor))
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	defer conn.Close()

	client = pb.NewStreamConversationsClient(conn)
	conversations()
}

func conversations() {
	stream, err := client.Converstaion(context.Background())
	if err != nil {
		log.Fatalf("connect failure", err)
	}
	for n := 0; n < 5; n++ {
		err := stream.Send(&pb.StreamRequest{Question: "stream client rpc " + strconv.Itoa(n)})
		if err != nil {
			log.Fatalf("stream request err %v ", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err %v", err)
		}
		log.Printf(res.Answer)
	}
}
