package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("port: ", os.Getenv("PORT"))
	svr := NewTCPServer(TCPServerConfig{
		Host: "localhost",
		Port: os.Getenv("PORT"),
	})
	svr.ListenAndServe()
}
