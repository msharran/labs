package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	pb "github.com/msharran/labs/go-grpc/pkg/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port       string
	serverAddr string
)

func main() {
	flag.StringVar(&port, "port", "3455", "http server port")
	flag.StringVar(&serverAddr, "server-addr", "localhost:3456", "gRPC server address")

	flag.Parse()

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	fmt.Println("Client connected to server at", serverAddr)

	// http server with /hello endpoint that calls the grpc server
	// http://localhost:8080/hello?name=Sharran

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s", reply.String())
	})

	fmt.Println("HTTP server listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
