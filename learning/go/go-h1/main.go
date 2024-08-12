package main

import (
	"context"
	"flag"
	"go-h1/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	port           string
	disableLogging bool
	disableAdmin   bool
)

func init() {
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&port, "p", "8080", "port to listen on (shorthand)")
}

func main() {
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx, cancel := context.WithCancel(context.Background())

	s := server.NewServer(server.ServerOpts{
		Ctx:  ctx,
		Addr: ":" + port,
		Log:  log,
	})

	// handle signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("interrupt signal received, shutting down")
		cancel()
	}()

	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	if err := s.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	s.Wait()
	log.Info("server exited gracefully")
}
