package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "go-http2/pkg/helloworld"

	"golang.org/x/net/http2"
)

func main() {
	// grpcServer := grpc.NewServer()
	// reflection.Register(grpcServer)
	// pb.RegisterGreeterServer(grpcServer, newServer())

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		v := r.PathValue("name")
		log.Printf("%s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello, %s!", v)
	})

	// use http2 server for grpc and http server

	http2Server := &http2.Server{}
	httpServer := &http.Server{
		Handler: mux,
	}
	err := http2.ConfigureServer(httpServer, http2Server)
	if err != nil {
		log.Fatalf("failed to configure http2 server: %v", err)
	}

	grpcServer := newServer()
	httpServer.ListenAndServe()

}

func newServer() *greeterServer {
	return &greeterServer{port: 3000}
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
