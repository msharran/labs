package main

import (
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	jl := slog.New(slog.NewJSONHandler(os.Stdout))
	slog.SetDefault(jl)
	slog.Info("Hello Gopher!",
		"from", "slog",
		"type", "json")

	tl := slog.New(slog.NewTextHandler(os.Stdout))
	slog.SetDefault(tl)
	slog.Info("Hello Gopher!",
		"from", "slog",
		"type", "text")
}
