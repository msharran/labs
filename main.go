package main

import "os"

func main() {
	svr := NewTCPServer(TCPServerConfig{
		Host: "localhost",
		Port: os.Getenv("PORT"),
	})
	svr.ListenAndServe()
}
