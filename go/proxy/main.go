package main

import (
	"flag"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/slog"
)

var serverAddr string
var timeout time.Duration

const proxyAddr = ":9090"

func init() {
	flag.StringVar(&serverAddr, "serverAddr", "localhost:9090", "Downstream address to proxy. Eg, `google.com:443`")
	t := flag.Int("timeout", 5, "Connection timeout in seconds. Default: 5")
	timeout = time.Duration(*t) * time.Second
}

func main() {
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout))
	slog.SetDefault(log)

	l, err := net.Listen("tcp", proxyAddr)
	check(err, "failed to setup TCP listener")

	slog.Info("started proxy", "Addr", proxyAddr)

	for {
		conn, err := l.Accept()
		check(err, "failed to accept connection")
		go proxy(conn)
	}
}

func proxy(upstreamConn net.Conn) {
	id := ulid.Make()
	log := slog.With("clientID", id)

	defer func() {
		upstreamConn.Close()
		log.Info("upstream connection closed")
	}()
	log.Info("connected to client")

	downstreamConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		slog.Error("error connecting to downstreamServer", err)
		return
	}
	defer func() {
		downstreamConn.Close()
		log.Info("downstream connection closed")
	}()

	// setting up conn deadlines
	upstreamConn.SetDeadline(time.Now().Add(timeout))
	downstreamConn.SetDeadline(time.Now().Add(timeout))

	slog.Info("sending traffic to downstreamAddr", "downstreamAddr", downstreamConn.RemoteAddr().String())

	var wg sync.WaitGroup
	wg.Add(2)

	// send traffic from upstreamConn to downstreamConn

	go func() {
		defer wg.Done()
		_, err := io.Copy(downstreamConn, upstreamConn)
		if err != nil {
			slog.Error("error copying upstream to downstream", err)
			return
		}
	}()

	// send traffic from downstreamConn to upstreamConn
	go func() {
		defer wg.Done()
		_, err = io.Copy(upstreamConn, downstreamConn)
		if err != nil {
			slog.Error("error copying downstream to upstream", err)
			return
		}
	}()

	wg.Wait()
}

func echoBack(conn net.Conn) {
	id := ulid.Make()
	log := slog.With("id", id)

	defer func() {
		conn.Close()
		log.Info("connection closed")
	}()

	log.Info("Connected to client")

	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		var buf [1024]byte
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
