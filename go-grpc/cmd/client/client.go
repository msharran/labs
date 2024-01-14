package main

import (
	"context"
	"flag"
	"fmt"

	pb "github.com/msharran/labs/go-grpc/pkg/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr string

func init() {
	flag.StringVar(&serverAddr, "server_addr", "localhost:3456", "gRPC server address")

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Sharran"})
	if err != nil {
		panic(err)
	}
	fmt.Println("reply:", reply.String())
}

func main() {
	flag.Parse()
}
