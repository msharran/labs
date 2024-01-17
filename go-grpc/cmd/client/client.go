package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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

	// http server with /hello endpoint that calls the grpc server
	// http://localhost:8080/hello?name=Sharran

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Error:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewGreeterClient(conn)
		fmt.Println("Client connected to server at", serverAddr)

		name := r.URL.Query().Get("name")
		reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			fmt.Println("Error:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		fmt.Fprintf(w, "%s", reply.String())
	})

	// /health endpoint to check if the server is up
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("HTTP server listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
