package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	pb "go-http2/pkg/helloworld"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port int

func init() {
	flag.IntVar(&port, "p", 3000, "The server port (shorthand)")
	flag.IntVar(&port, "port", 3000, "The server port")
}

func main() {
	flag.Parse()

	// configure http server
	h1Mux := http.NewServeMux()
	h1Mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		v := r.PathValue("name")
		log.Printf("%s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello, %s!", v)
	})

	// configure gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterGreeterServer(grpcServer, configureGRPCServer())

	// configure h2c server
	http2Server := http2.Server{}
	h2cHttpServer := http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				h1Mux.ServeHTTP(w, r)
			}
		}), &http2Server),
	}

	err := http2.ConfigureServer(&h2cHttpServer, &http2Server)
	if err != nil {
		log.Fatal("unable set up serve", err)
	}

	// Listen
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Fatal(h2cHttpServer.Serve(lis))
}

func configureGRPCServer() *greeterServer {
	return &greeterServer{}
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements gogrpc.GreeterServer
func (*greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", r.GetName())
	return &pb.HelloReply{Message: "Hello, " + r.Name}, nil
}
