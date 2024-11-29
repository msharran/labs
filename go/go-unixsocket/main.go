package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const SOCK = "/tmp/foo.sock"

func main() {
	socket, err := net.Listen("unix", SOCK)
	check(err)

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("removing %s\n", SOCK)
		os.Remove(SOCK)
		os.Exit(1)
	}()

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello curious developer!"))
	})

	server := http.Server{
		Handler: m,
	}

	fmt.Println("starting server")
	err = server.Serve(socket)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
