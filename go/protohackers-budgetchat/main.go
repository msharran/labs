package main

import (
	"budgetchat/internal/chatserver"
	"flag"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

var timeout time.Duration

const serverAddr = ":10000"

func init() {
	t := flag.Int("timeout", 5, "Connection timeout in seconds. Default: 5")
	timeout = time.Duration(*t) * time.Second
}

func main() {
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout))

	svr := chatserver.NewServer(chatserver.ServerArgs{
		Log:     log,
		Timeout: timeout,
	})

	if err := svr.Listen(serverAddr); err != nil {
		log.Error("server failed", err)
	}
}
