package main

import (
	"net"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/slog"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout))
	slog.SetDefault(log)
	slog.Info("setting up proxy")

	l, err := net.Listen("tcp", ":8080")
	check(err, "failed to setup TCP listener")

	for {
		conn, err := l.Accept()
		check(err, "failed to accept connection")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	id := ulid.Make()
	log := slog.With("id", id)

	defer func() {
		conn.Close()
		log.Info("connection closed")
	}()

	log.Info("Connected to client")

	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		var buf [128]byte
		read, err := conn.Read(buf[:])
		if err != nil {
			slog.Error("failed to read from connection", err, "bytes", read)
			return
		}

		wrote, err := os.Stderr.Write(buf[:read])
		if err != nil {
			slog.Error("failed to write to stderr", err, "bytes", wrote)
			return
		}
		slog.Info("wrote to stderr", "read", read, "wrote", wrote)
	}
}

func check(err error, msg string) {
	if err != nil {
		slog.Error(msg, err)
		os.Exit(1)
	}
}
