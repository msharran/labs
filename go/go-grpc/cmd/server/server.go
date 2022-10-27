package main

import (
	context "context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/msharran/labs/go-grpc/internal/helloworld"
	"google.golang.org/grpc"
)

var port int

func init() {
	flag.IntVar(&port, "port", 3456, "port on with the gRPC server should serve connections")
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, newServer())

	fmt.Printf("gRPC: serving on %s\n", lis.Addr().String())
	grpcServer.Serve(lis)
}

func newServer() *greeterServer {
	return &greeterServer{port: port}
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
	port int
}

// SayHello implements gogrpc.GreeterServer
func (*greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("request:", r.String())
	return &pb.HelloReply{Message: "Hello, " + r.Name}, nil
}
