package main

import (
	"bufio"
	"net"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

var (
	Host = "localhost"
	Port = "8080"
)

func main() {
	// lgr := slog.New(slog.NewTextHandler(os.Stdout))
}

type Server interface {
	Serve(l net.Listener)
}

type server struct {
	log  *slog.Logger
	host string
	port string
}

func (s *server) Serve(l net.Listener) {
	s.log.Info("accepting incoming connections")
	for {
		conn, err := l.Accept()
		log := s.log.With("addr", conn.LocalAddr().String())
		if err != nil {
			log.Error("error accepting connection ", err)
			os.Exit(1)
		}
		go s.handleRequest(conn)
	}
}

func (s *server) handleRequest(conn net.Conn) {
	log := s.log.With("request-id", uuid.NewString(), "addr", conn.LocalAddr().String())

	log.Info("handling connection")

	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("unable to close the connection", nil)
		}
	}()

	r := bufio.NewReader(conn)
	for {
		log.Info("processing connection")
		msg := make([]byte, 9)
		_, err := r.Read(msg)
		if err != nil {
			log.Error("error reading from connection", err)
			break
		}
	}
}
